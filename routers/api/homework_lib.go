package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"net/http"
	"strconv"
)

type CreateHomeworkLibForm struct {
	CourseId uint   `form:"courseId" binding:"required"`
	Name     string `form:"name" binding:"required"`
}

func CreateHomeworkLib(c *gin.Context) {
	var form CreateHomeworkLibForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibService{
		CourseId: form.CourseId,
		Name:     form.Name,
	}
	if _, err := s.AddHomeworkLib(); err == nil {
		c.String(http.StatusCreated, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type UpdateHomeworkLibForm struct {
	Id   uint   `form:"id" binding:"required"`
	Name string `form:"name"`
}

func UpdateHomeworkLibNameById(c *gin.Context) {
	var form UpdateHomeworkLibForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibService{
		Id:   form.Id,
		Name: form.Name,
	}
	if err := s.UpdateHomeworkLibById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

func GetHomeworkLibsByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibService{
		CourseId: uint(courseId),
	}
	if libs, err := s.GetHomeworkLibsByCourseId(); err == nil {
		c.JSON(http.StatusOK, libs)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetHomeworkLibWithItemsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibService{
		Id: uint(id),
	}
	if lib, err := s.GetHomeworkLibWithItemsById(); err == nil {
		for _, v := range lib.Items { //获取作业题目信息时不把正确答案返回
			v.Answer = ""
		}
		c.JSON(http.StatusOK, lib)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
