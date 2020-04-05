package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"net/http"
)

type CreateClassForm struct {
	CourseId  uint   `form:"courseId" binding:"required" `
	ClassName string `form:"className" binding:"required" `
}

func CreateClass(c *gin.Context) {
	var form CreateClassForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
	} else {
		if claims, ok := c.Get("claims"); ok {
			s := service.ClassService{
				CourseId:  form.CourseId,
				UserId:    claims.(*utils.Claims).Id,
				ClassName: form.ClassName,
			}
			if _, err := s.CreateClass(); err != nil {
				c.String(http.StatusInternalServerError, "")
			} else {
				c.String(http.StatusCreated, "")
			}
		}
	}

}
