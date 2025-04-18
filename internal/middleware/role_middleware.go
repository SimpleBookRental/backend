package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/models"
)

// RequireRole middleware checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Unauthorized: role not found in context",
				"data":    nil,
			})
			return
		}

		// Check if user has the required role
		if role != requiredRole {
			c.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"message": "Forbidden: insufficient permissions",
				"data":    nil,
			})
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware checks if the user is an admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.AdminRole)
}

// RequireUser middleware checks if the user is a regular user
func RequireUser() gin.HandlerFunc {
	return RequireRole(models.UserRole)
}

// RequireAdminOrSameUser middleware checks if the user is an admin or the same user
func RequireAdminOrSameUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role and user ID from context
		role, roleExists := c.Get("role")
		userID, userIDExists := c.Get("user_id")

		// Check if role and user ID exist in context
		if !roleExists || !userIDExists {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Unauthorized: user information not found in context",
				"data":    nil,
			})
			return
		}

		// If user is admin, allow access
		if role == models.AdminRole {
			c.Next()
			return
		}

		// Get the requested user ID from the URL parameter
		requestedUserID := c.Param("id")

		// If no user ID in URL or it's the same user, allow access
		if requestedUserID == "" || requestedUserID == userID {
			c.Next()
			return
		}

		// Otherwise, deny access
		c.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"message": "Forbidden: insufficient permissions",
			"data":    nil,
		})
	}
}
