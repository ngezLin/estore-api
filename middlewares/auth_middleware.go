package middlewares

import (
	"estore-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewareAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		//check auth
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization is required"})
			c.Abort()
			return
		}

		//baca auth header
		parts := strings.Split(authHeader, " ")

		//baca index pertama
		AdminId, err := utils.ValidateTokenAdmin(parts[1])

		if err != nil {
			c.JSON(401, gin.H{"error": "invalid Token"})
			c.Abort()
			return
		}

		//setting key adminId dan value adminId
		c.Set("adminId", AdminId)
		c.Next()
	}
}
