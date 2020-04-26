package api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"learning/conf"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

type CreateChapterForm struct {
	CourseId    uint   `json:"courseId" binding:"required" `
	ChapterName string `json:"chapterName" binding:"required" `
}

func CreateChapter(c *gin.Context) {
	var form CreateChapterForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.ChapterService{
			UserId:      claims.(*utils.Claims).Id,
			CourseId:    form.CourseId,
			ChapterName: form.ChapterName,
		}
		if _, err := s.AddChapter(); err == nil {
			c.String(http.StatusCreated, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
func GetChapterByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ChapterService{
		CourseId: uint(courseId),
	}
	if chapters, err := s.GetChapterByCourseId(); err == nil {
		c.JSON(http.StatusOK, chapters)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type UpdateChapterNameForm struct {
	Id          uint   `json:"courseId" binding:"required" `
	ChapterName string `json:"chapterName" binding:"required" `
}

func UpdateChapterName(c *gin.Context) {
	var form UpdateChapterNameForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	s := service.ChapterService{
		Id:          form.Id,
		ChapterName: form.ChapterName,
	}
	if err := s.UpdateChapterById(); err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}
func DeleteChapterById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ChapterService{
		Id: uint(id),
	}
	chapter, _ := s.GetChapterById()
	if chapter == nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if err = s.DeleteChapterById(); err == nil {
		if chapter.VideoName != nil && *chapter.VideoName != "" {
			filepath := conf.AppConfig.Path.Video + "/" + *chapter.VideoName
			if err = os.Remove(filepath); err != nil {
				log.Println("视频删除失败", err)
			}
		}
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}

}

type UpdateChapterVideoForm struct {
	Id uint `json:"id" binding:"required" `
}

func UpdateChapterVideo(c *gin.Context) {
	var form UpdateChapterVideoForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	file, err := c.FormFile("video")
	if err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	u1 := uuid.Must(uuid.NewV4(), nil).String()
	filepath := conf.AppConfig.Path.Video + "/" + u1 + ".mp4"
	if err := c.SaveUploadedFile(file, filepath); err == nil {
		videoName := u1 + ".mp4"
		s := service.ChapterService{
			Id:        form.Id,
			VideoName: &videoName,
		}
		if err = s.UpdateChapterById(); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"videoName": videoName,
			})
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
func DeleteChapterVideoById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	var empty string
	s := service.ChapterService{
		Id:        uint(id),
		VideoName: &empty,
	}
	chapter, _ := s.GetChapterById()
	if chapter == nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if err = s.UpdateChapterById(); err == nil {
		if chapter.VideoName != nil && *chapter.VideoName != "" {
			filepath := conf.AppConfig.Path.Video + "/" + *chapter.VideoName
			if err = os.Remove(filepath); err != nil {
				log.Println("视频删除失败", err)
			}
		}
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
