package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"net/http"
	"strconv"
)

type CreateExamLibForm struct {
	CourseId uint   `form:"courseId" binding:"required"`
	Name     string `form:"name" binding:"required"`
}

func CreateExamLib(c *gin.Context) {
	var form CreateExamLibForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibService{
		CourseId: form.CourseId,
		Name:     form.Name,
	}
	if _, err := s.AddExamLib(); err == nil {
		c.String(http.StatusCreated, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type UpdateExamLibForm struct {
	Id   uint   `form:"id" binding:"required"`
	Name string `form:"name"`
}

func UpdateExamLibNameById(c *gin.Context) {
	var form UpdateExamLibForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibService{
		Id:   form.Id,
		Name: form.Name,
	}
	if err := s.UpdateExamLibById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamLibsByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibService{
		CourseId: uint(courseId),
	}
	if libs, err := s.GetExamLibsByCourseId(); err == nil {
		c.JSON(http.StatusOK, libs)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamLibWithItemsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibService{
		Id: uint(id),
	}
	if lib, err := s.GetExamLibWithItemsById(); err == nil {
		if c.Query("answer") == "" {
			for _, v := range lib.Items { //获取作业题目信息时不把正确答案返回
				v.Answer = ""
			}
		}
		c.JSON(http.StatusOK, lib)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
