package api

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimiter tracks request rates by IP
type RateLimiter struct {
	mu         sync.Mutex
	limits     map[string]*IPLimit
	rate       int
	window     time.Duration
	lastClean  time.Time
	cleanEvery time.Duration
}

// IPLimit tracks requests for a single IP
type IPLimit struct {
	count    int
	resetAt  time.Time
}

// Middleware holds all middleware handlers
type Middleware struct {
	jwtService  *auth.JWTService
	logger      *logger.Logger
	rateLimiter *RateLimiter
}

// NewMiddleware creates a new Middleware
func NewMiddleware(jwtService *auth.JWTService, logger *logger.Logger) *Middleware {
	// Create rate limiter: 100 requests per minute
	rateLimiter := &RateLimiter{
		limits:     make(map[string]*IPLimit),
		rate:       100,
		window:     time.Minute,
		lastClean:  time.Now(),
		cleanEvery: 5 * time.Minute,
	}

	return &Middleware{
		jwtService:  jwtService,
		logger:      logger,
		rateLimiter: rateLimiter,
	}
}

// LoggerMiddleware logs each request
func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		
		m.logger.Info("API Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// RecoveryMiddleware recovers from panics
func (m *Middleware) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("Recovered from panic", zap.Any("error", err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// CORSMiddleware sets CORS headers
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthMiddleware checks if the user is authenticated
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		// Extract token
		token := authHeader[7:] // Remove "Bearer " prefix

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Error("Invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Token is valid, add user ID and role to context
		c.Set("userId", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if the user has the required role
func (m *Middleware) RoleMiddleware(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check role
		switch role {
		case domain.RoleAdmin:
			if userRole != string(domain.RoleAdmin) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
				return
			}
		case domain.RoleLibrarian:
			if userRole != string(domain.RoleAdmin) && userRole != string(domain.RoleLibrarian) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Librarian access required"})
				return
			}
		case domain.RoleMember:
			if userRole != string(domain.RoleAdmin) && userRole != string(domain.RoleLibrarian) && userRole != string(domain.RoleMember) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Member access required"})
				return
			}
		}

		c.Next()
	}
}

// RateLimitMiddleware limits the number of requests
func (m *Middleware) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		// Get current limits for this IP
		reached, remaining, resetAt := m.checkRateLimit(ip)
		
		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(m.rateLimiter.rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", resetAt.Format(time.RFC3339))
		
		// Check if limit is exceeded
		if reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}
		
		c.Next()
	}
}

// checkRateLimit verifies if the rate limit is reached and returns status
func (m *Middleware) checkRateLimit(ip string) (bool, int, time.Time) {
	m.rateLimiter.mu.Lock()
	defer m.rateLimiter.mu.Unlock()
	
	now := time.Now()
	
	// Periodically clean old entries
	if now.Sub(m.rateLimiter.lastClean) > m.rateLimiter.cleanEvery {
		for key, limit := range m.rateLimiter.limits {
			if now.After(limit.resetAt) {
				delete(m.rateLimiter.limits, key)
			}
		}
		m.rateLimiter.lastClean = now
	}
	
	// Get or create limit for this IP
	limit, exists := m.rateLimiter.limits[ip]
	if !exists || now.After(limit.resetAt) {
		// Create new limit
		limit = &IPLimit{
			count:    0,
			resetAt:  now.Add(m.rateLimiter.window),
		}
		m.rateLimiter.limits[ip] = limit
	}
	
	// Increment count
	limit.count++
	
	// Check if over limit
	reached := limit.count > m.rateLimiter.rate
	remaining := m.rateLimiter.rate - limit.count
	if remaining < 0 {
		remaining = 0
	}
	
	return reached, remaining, limit.resetAt
}
