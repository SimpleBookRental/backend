package middleware

import (
	"strings"

	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
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
