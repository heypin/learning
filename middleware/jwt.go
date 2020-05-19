package middleware

import (
	"github.com/gin-gonic/gin"
	"learning/utils"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			if claims, err := utils.ParseToken(token); err == nil {
				c.Set("claims", claims)
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, "")
		c.Abort()
	}
}

//path := c.Request.URL.Path
//if (claims.Role == models.RoleAdmin && strings.HasPrefix("/admin/", path)) ||
//(claims.Role == models.RoleUser) {
//c.Set("claims", claims)
//c.Next()
//return
//}
