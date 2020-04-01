package routers

import (
	"github.com/gin-gonic/gin"
	"learning/conf"
	"learning/middleware"
	"learning/routers/api"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Static("/index.html", conf.AppConfig.Path.Frontend)
	r.Static("/avatar", conf.AppConfig.Path.Avatar)
	r.Static("/file", conf.AppConfig.Path.File)
	r.POST("/login", api.UserLogin)
	r.POST("/register", api.UserRegister)
	auth := r.Group("/")
	auth.Use(middleware.JWT())
	{
		auth.GET("user", api.GetUserByToken)
		auth.PUT("user", api.UpdateUserById)
		auth.PUT("user/password", api.UpdateUserPassword)
	}
	return r
}
