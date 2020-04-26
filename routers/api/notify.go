package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"strconv"
)

type CreateNotifyForm struct {
	CourseId uint   `json:"courseId" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func CreateNotify(c *gin.Context) {
	var form CreateNotifyForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.NotifyService{
			UserId:   claims.(*utils.Claims).Id,
			CourseId: form.CourseId,
			Title:    form.Title,
			Content:  form.Content,
		}
		if _, err := s.AddNotify(); err == nil {
			c.String(http.StatusCreated, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

type UpdateNotifyForm struct {
	Id      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func UpdateNotifyById(c *gin.Context) {
	var form UpdateNotifyForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	s := service.NotifyService{
		Id:      form.Id,
		Title:   form.Title,
		Content: form.Content,
	}
	if err := s.UpdateNotifyById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

func GetNotifyByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.NotifyService{
		CourseId: uint(courseId),
	}
	if notifies, err := s.GetNotifyByCourseId(); err == nil {
		c.JSON(http.StatusOK, notifies)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func DeleteNotifyById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.NotifyService{
		Id: uint(id),
	}
	if err = s.DeleteNotifyById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
