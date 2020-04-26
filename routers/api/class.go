package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"net/http"
	"strconv"
)

type CreateClassForm struct {
	CourseId  uint   `json:"courseId" binding:"required" `
	ClassName string `json:"className" binding:"required" `
}

func CreateClass(c *gin.Context) {
	var form CreateClassForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ClassService{
		CourseId:  form.CourseId,
		ClassName: form.ClassName,
	}
	if _, err := s.CreateClass(); err == nil {
		c.String(http.StatusCreated, "")
		return
	}
	c.String(http.StatusInternalServerError, "")

}
func GetClassByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ClassService{
		CourseId: uint(courseId),
	}
	if classes, err := s.GetClassByCourseId(); err == nil {
		c.JSON(http.StatusOK, classes)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
