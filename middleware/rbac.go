package middleware

import (
	"net/http"

	"andrei-api/models"
	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		currentUser := user.(models.User)
		
		for _, role := range allowedRoles {
			if currentUser.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	})
}

func RequireAndrei() gin.HandlerFunc {
	return RequireRole(models.RoleAndrei)
}

func RequireDemon() gin.HandlerFunc {
	return RequireRole(models.RoleDemon)
}

func RequireNetworkAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleNetworkAdmin)
}

func RequireAndreiOrDemon() gin.HandlerFunc {
	return RequireRole(models.RoleAndrei, models.RoleDemon)
}