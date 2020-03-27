package api

import (
	"github.com/gin-gonic/gin"
	"learning/models"
	"learning/service"
	"learning/utils"
	"net/http"
)

type UserLoginForm struct {
	Email    string `form:"email" binding:"required,email" `
	Password string `form:"password" binding:"required,min=8" `
	Role     int    `form:"role" binding:"required" `
}

func UserLogin(c *gin.Context) {
	var user UserLoginForm
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
	} else {
		if user.Role == models.ROLE_STUDENT {
			s := service.StudentService{
				Email:    user.Email,
				Password: user.Password,
			}
			if ok := s.Auth(); ok {
				token, _ := utils.GenerateToken(s.Email, models.ROLE_STUDENT)
				c.JSON(http.StatusOK, gin.H{
					"token": token,
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"err": "帐号或密码错误"})
			}
		} else if user.Role == models.ROLE_TEACHER {
			t := service.TeacherService{
				Email:    user.Email,
				Password: user.Password,
			}
			if ok := t.Auth(); ok {
				token, _ := utils.GenerateToken(t.Email, models.ROLE_TEACHER)
				c.JSON(http.StatusOK, gin.H{
					"token": token,
				})
			}
		}

	}
}
func UserRegister(c *gin.Context) {
	var s service.StudentService
	if err := c.ShouldBind(&s); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		if _, err := s.Register(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "该邮箱已被注册",
			})
		} else {
			c.JSON(http.StatusCreated, "")
		}
	}
}
