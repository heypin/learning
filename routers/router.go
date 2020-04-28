package routers

import (
	"github.com/gin-gonic/gin"
	"learning/conf"
	"learning/middleware"
	"learning/routers/api"
)

func InitRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Static("/learning", conf.AppConfig.Path.Frontend)
	r.Static("/avatar", conf.AppConfig.Path.Avatar)
	r.Static("/resource", conf.AppConfig.Path.File)
	r.Static("/cover", conf.AppConfig.Path.Cover)

	r.POST("/login", api.UserLogin)
	r.POST("/register", api.UserRegister)
	r.GET("register/code", api.GenerateRegisterCode)
	r.GET("/video/:name", api.PlayVideo)
	r.POST("/compile", api.ExecuteProgram)
	r.MaxMultipartMemory = 500 << 20 //500MB

	auth := r.Group("/")
	auth.Use(middleware.JWT())
	{
		auth.GET("user", api.GetUserByToken)
		auth.PUT("user", api.UpdateUserById)
		auth.PUT("user/password", api.UpdateUserPassword)

		auth.POST("course", api.AddCourse)
		auth.GET("course/teach", api.GetUserTeachCourse)

		auth.GET("class", api.GetClassByCourseId)
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

		auth.GET("notify", api.GetNotifyByCourseId)
		auth.POST("notify", api.CreateNotify)
		auth.PUT("notify", api.UpdateNotifyById)
		auth.DELETE("notify", api.DeleteNotifyById)

		auth.GET("comment", api.GetCommentByCourseId)
		auth.GET("comment/user", api.GetCommentByUserId)
		auth.GET("comment/reply", api.GetCommentReplyToUser)
		auth.POST("comment", api.CreateComment)
		auth.DELETE("comment", api.DeleteCommentById)

		auth.GET("classMember/user", api.GetUsersByClassId)
		auth.GET("classMember/class", api.GetClassesByUserId)
		auth.POST("classMember/join", api.JoinClassByClassCode)
		auth.DELETE("classMember", api.DeleteClassMember)

		auth.GET("homeworkLib", api.GetHomeworkLibsByCourseId)
		auth.GET("homeworkLib/items", api.GetHomeworkLibWithItemsById)
		auth.PUT("homeworkLib/name", api.UpdateHomeworkLibNameById)
		auth.POST("homeworkLib", api.CreateHomeworkLib)

		auth.GET("homeworkLibItem", api.GetHomeworkLibItemsByLibId)
		auth.PUT("homeworkLibItem", api.UpdateHomeworkLibItemAndOptions)
		auth.POST("homeworkLibItem", api.CreateHomeworkLibItemAndOptions)
		auth.DELETE("homeworkLibItem", api.DeleteHomeworkLibItemById)

		auth.GET("homeworkPublish", api.GetHomeworkPublishById)
		auth.GET("homeworkPublish/class", api.GetHomeworkPublishesByClassId)
		auth.GET("homeworkPublish/submit", api.GetHomeworkPublishesWithSubmitByClassId)
		auth.POST("homeworkPublish", api.PublishHomework)
		auth.PUT("homeworkPublish", api.UpdateHomeworkPublishById)

		auth.GET("homeworkSubmit", api.GetHomeworkSubmitById)
		auth.GET("homeworkSubmit/publish", api.GetHomeworkSubmitsByPublishId)
		auth.GET("homeworkSubmit/user", api.GetHomeworkUserSubmitWithItems)
		auth.POST("homeworkSubmit", api.SubmitHomeworkWithItems)
		auth.PUT("homeworkSubmit/mark", api.UpdateHomeworkSubmitMarkById)
		auth.PUT("homeworkSubmit/score", api.UpdateHomeworkSubmitItemsScore)

		auth.GET("examLib", api.GetExamLibsByCourseId)
		auth.GET("examLib/items", api.GetExamLibWithItemsById)
		auth.PUT("examLib/name", api.UpdateExamLibNameById)
		auth.POST("examLib", api.CreateExamLib)

		auth.GET("examLibItem", api.GetExamLibItemsByLibId)
		auth.PUT("examLibItem", api.UpdateExamLibItemAndOptions)
		auth.POST("examLibItem", api.CreateExamLibItemAndOptions)
		auth.DELETE("examLibItem", api.DeleteExamLibItemById)

		auth.GET("examPublish", api.GetExamPublishById)
		auth.GET("examPublish/class", api.GetExamPublishesByClassId)
		auth.GET("examPublish/submit", api.GetExamPublishesWithSubmitByClassId)
		auth.POST("examPublish", api.PublishExam)
		auth.PUT("examPublish", api.UpdateExamPublishById)

		auth.GET("examSubmit", api.GetExamSubmitById)
		auth.GET("examSubmit/publish", api.GetExamSubmitsByPublishId)
		auth.GET("examSubmit/user", api.GetExamUserSubmitWithItems)
		auth.POST("examSubmit/item", api.SubmitExamItem)
		auth.POST("examSubmit/start", api.StartExam)
		auth.PUT("examSubmit/finish", api.FinishExam)
		auth.PUT("examSubmit/score", api.UpdateExamSubmitItemsScore)

		auth.GET("exam/excel", api.ExportExamToExcel)
		auth.GET("homework/excel", api.ExportHomeworkToExcel)
	}
	return r
}
