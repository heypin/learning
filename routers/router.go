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
	r.Static("/learning", conf.AppConfig.Path.Frontend)
	r.Static("/avatar", conf.AppConfig.Path.Avatar)
	r.Static("/resource", conf.AppConfig.Path.File)
	r.Static("/cover", conf.AppConfig.Path.Cover)
	r.POST("/login", api.UserLogin)
	r.POST("/register", api.UserRegister)
	r.GET("/video/:name", api.PlayVideo)
	r.MaxMultipartMemory = 500 << 20 //500MB
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

		auth.GET("chapter", api.GetChapterByCourseId)
		auth.POST("chapter", api.CreateChapter)
		auth.PUT("chapter", api.UpdateChapterName)
		auth.DELETE("chapter", api.DeleteChapterById)
		auth.DELETE("chapter/video", api.DeleteChapterVideoById)
		auth.PUT("chapter/video", api.UpdateChapterVideo)

	}
	return r
}
