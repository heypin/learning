package routers

import (
	"github.com/gin-gonic/gin"
	"learning/middleware"
	"learning/routers/api"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/test", api.Test)
	r.POST("/login", api.UserLogin)
	auth := r.Group("/")
	auth.Use(middleware.JWT())
	{
		auth.GET("student/", api.GetStudentByToken)
	}
	return r
}
