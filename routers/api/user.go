package api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"learning/cache"
	"learning/conf"
	"learning/models"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
)

type UserLoginForm struct {
	Email    string `json:"email" binding:"required,email" `
	Password string `json:"password" binding:"required,min=8" `
}

func UserLogin(c *gin.Context) {
	var form UserLoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
	} else {
		s := service.UserService{
			Email:    form.Email,
			Password: form.Password,
		}
		if id, ok := s.Auth(); ok {
			token, _ := utils.GenerateToken(id, models.RoleUser)
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"err": "帐号或密码错误"})
		}
	}
}

type UserRegisterForm struct {
	Email    string `json:"email" binding:"required,email" `
	Password string `json:"password" binding:"required,min=8" `
	RealName string `json:"realName" binding:"required" `
	Captcha  string `json:"captcha" binding:"required"`
}

func UserRegister(c *gin.Context) {
	var form UserRegisterForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "参数错误",
		})
		return
	}
	if value, err := cache.RedisClient.Get(cache.CaptchaPrefix + "." + form.Email).Result(); err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	} else {
		if value != form.Captcha {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "验证码错误",
			})
			return
		}
	}
	s := service.UserService{
		Email:    form.Email,
		Password: form.Password,
		RealName: form.RealName,
	}
	if _, err := s.Register(); err == nil {
		c.String(http.StatusCreated, "")
	} else if err == models.ErrRecordExist {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "该邮箱已被注册",
		})
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetUserByToken(c *gin.Context) {
	if claims, ok := c.Get("claims"); ok {
		s := service.UserService{
			Id: claims.(*utils.Claims).Id,
		}
		user, err := s.GetUserById()
		if err == nil && user != nil {
			c.JSON(http.StatusOK, gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"avatar":   user.Avatar,
				"realName": user.RealName,
				"sex":      user.Sex,
				"number":   user.Number,
			})
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

type UserUpdateForm struct {
	Email    string `form:"email" binding:"required,email" `
	RealName string `form:"realName" binding:"required" `
	Number   string `form:"number"`
	Sex      uint   `form:"sex" `
}

func UpdateUserById(c *gin.Context) {
	var form UserUpdateForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		if claims, ok := c.Get("claims"); ok {
			var avatar string = ""
			if file, err := c.FormFile("avatar"); err == nil {
				u1 := uuid.Must(uuid.NewV4(), nil).String()
				filepath := conf.AppConfig.Path.Avatar + "/" + u1 + ".png"
				if err := c.SaveUploadedFile(file, filepath); err == nil {
					avatar = u1 + ".png"
					log.Println("上传头像成功")
				}
			} else {
				log.Println("上传头像失败")
			}
			s := service.UserService{
				Id:       claims.(*utils.Claims).Id,
				Email:    form.Email,
				RealName: form.RealName,
				Number:   form.Number,
				Avatar:   avatar,
				Sex:      form.Sex,
			}
			if err := s.UpdateUserById(); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"avatar": avatar,
				})
				return
			}
		}
		c.String(http.StatusInternalServerError, "")
	}
}

type UserPasswordForm struct {
	Password    string `json:"password" binding:"required,min=8" `
	OldPassword string `json:"oldPassword" binding:"required,min=8" `
}

func UpdateUserPassword(c *gin.Context) {
	var form UserPasswordForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.UserService{
			Id:       claims.(*utils.Claims).Id,
			Password: form.Password,
		}
		if err := s.UpdateUserPassword(form.OldPassword); err == nil {
			c.String(http.StatusOK, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")

}
