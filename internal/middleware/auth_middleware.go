package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(tokenRepo *repositories.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Authorization header is required",
				"data":    nil,
			})
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Authorization header format must be Bearer {token}",
				"data":    nil,
			})
			return
		}

		// Get token
		tokenString := parts[1]

		// Check if token exists and is valid
		issuedToken, err := tokenRepo.FindTokenByValue(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Error validating token",
				"data":    nil,
			})
			return
		}

		// Check if token exists
		if issuedToken == nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid token",
				"data":    nil,
			})
			return
		}

		// Check if token is revoked
		if issuedToken.IsRevoked {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token has been revoked",
				"data":    nil,
			})
			return
		}

		// Check if token is expired
		if issuedToken.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token has expired",
				"data":    nil,
			})
			return
		}

		// Check if token is an access token
		if issuedToken.TokenType != string(models.AccessToken) {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid token type",
				"data":    nil,
			})
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, utils.GetAccessSecret())
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"data":    nil,
			})
			return
		}

		// Set user ID and email in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
