package middleware

import (
	"github.com/gin-gonic/gin"
	"learning/models"
	"learning/utils"
	"net/http"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			claims, err := utils.ParseToken(token)
			if err == nil {
				path := c.Request.URL.Path
				if (claims.Role == models.ROLE_ADMIN && strings.HasPrefix("/admin/", path)) ||
					(claims.Role == models.ROLE_USER) {
					c.Set("claims", claims)
					c.Next()
					return
				}
			}
		}
		c.JSON(http.StatusUnauthorized, "")
		c.Abort()
	}
}
