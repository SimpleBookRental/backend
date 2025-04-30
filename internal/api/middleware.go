package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Middleware represents API middleware
type Middleware struct {
	jwtService *auth.JWTService
	logger     *logger.Logger
}

// NewMiddleware creates a new middleware
func NewMiddleware(jwtService *auth.JWTService, logger *logger.Logger) *Middleware {
	return &Middleware{
		jwtService: jwtService,
		logger:     logger.Named("middleware"),
	}
}

// AuthMiddleware is a middleware for authentication
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		// Extract token from header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Error("Invalid token", zap.Error(err))
			SendError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		// Check token type
		if claims.TokenType != auth.AccessToken {
			m.logger.Error("Token is not an access token", zap.String("tokenType", string(claims.TokenType)))
			SendError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user ID and role in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RoleMiddleware is a middleware for role-based access control
func (m *Middleware) RoleMiddleware(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		roleInterface, exists := c.Get("userRole")
		if !exists {
			SendError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		role := domain.UserRole(roleInterface.(string))

		// Check if user has required role
		hasRole := false
		for _, r := range roles {
			if role == r {
				hasRole = true
				break
			}
		}

		// Admin has access to everything
		if role == domain.RoleAdmin {
			hasRole = true
		}

		if !hasRole {
			SendError(c, domain.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware is a middleware for CORS
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware is a middleware for logging
func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		m.logger.Info("Request",
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.Int("status_code", statusCode),
			zap.String("path", path),
			zap.Duration("latency", latency),
		)
	}
}

// RecoveryMiddleware is a middleware for recovering from panics
func (m *Middleware) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("Panic recovered", zap.Any("error", err))
				SendError(c, domain.ErrInternalServer)
				c.Abort()
			}
		}()
		c.Next()
	}
}
