package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
)

func GetStudentByToken(c *gin.Context) {
	if claims, ok := c.Get("claims"); ok {
		log.Println(claims.(*utils.Claims).Email)
		s := service.StudentService{
			Email: claims.(*utils.Claims).Email,
		}
		student, err := s.GetStudentByEmail()
		if err == nil && student != nil {
			c.JSON(http.StatusOK, gin.H{
				"id":       student.ID,
				"email":    student.Email,
				"avatar":   student.Avatar,
				"realName": student.RealName,
			})
		}
	}
}

func Test(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{})
	c.JSON(http.StatusOK, nil)
}
