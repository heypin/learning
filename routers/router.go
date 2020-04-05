package routers

import (
	"fmt"
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
	r.Static("/resource", conf.AppConfig.Path.File)
	r.Static("/cover", conf.AppConfig.Path.Cover)
	r.POST("/login", api.UserLogin)
	r.POST("/register", api.UserRegister)
	r.MaxMultipartMemory = 500 << 20 //500MB
	r.GET("/download", func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "123.png"))
		c.File(conf.AppConfig.Path.File + "/" + "2d38add9-802f-4821-8498-1ccfd75aa2b7.png")
	})
	auth := r.Group("/")
	auth.Use(middleware.JWT())
	{
		auth.GET("user", api.GetUserByToken)
		auth.PUT("user", api.UpdateUserById)
		auth.PUT("user/password", api.UpdateUserPassword)
		auth.POST("course", api.AddCourse)
		auth.GET("course/teach", api.GetTeachCourse)
		auth.POST("class", api.CreateClass)
		auth.GET("file/children", api.GetChildFile)
		auth.GET("file/download", api.DownloadFile)
		auth.POST("file", api.CreateFile)
		auth.POST("file/folder", api.CreateFolder)
		auth.DELETE("file", api.DeleteFile)
	}
	return r
}
