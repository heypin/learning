package middleware

import (
	"github.com/gin-gonic/gin"
	"learning/models"
	"learning/utils"
	"log"
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
				log.Println("jwt ", path)
				if (claims.Role == models.ROLE_ADMIN && strings.HasPrefix("/admin/", path)) ||
					(claims.Role == models.ROLE_STUDENT && strings.HasPrefix("/student/", path)) ||
					(claims.Role == models.ROLE_TEACHER && strings.HasPrefix("/teacher/", path)) {
					log.Println(claims)
					c.Set("claims", claims)
					c.Next()
				}
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
		}

	}
}
