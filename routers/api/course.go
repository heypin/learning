package api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"learning/conf"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
)

type CreateCourseForm struct {
	Name        string `form:"name" binding:"required" `
	Teacher     string `form:"teacher" binding:"required" `
	Description string `form:"description" `
}

func AddCourse(c *gin.Context) {
	var form CreateCourseForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
	} else {
		if claims, ok := c.Get("claims"); ok {
			var cover string = ""
			if file, err := c.FormFile("cover"); err == nil {
				u1 := uuid.Must(uuid.NewV4(), nil).String()
				filepath := conf.AppConfig.Path.Cover + "/" + u1 + ".jpg"
				if err := c.SaveUploadedFile(file, filepath); err == nil {
					cover = u1 + ".jpg"
					log.Println("上传封面成功")
				} else {
					log.Println("上传封面失败")
				}
			} else {
				log.Println("未获取到封面文件")
			}
			s := service.CourseService{
				UserId:      claims.(*utils.Claims).Id,
				Name:        form.Name,
				Teacher:     form.Teacher,
				Cover:       cover,
				Description: form.Description,
			}
			if _, err := s.AddCourse(); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"cover": cover,
				})
				return
			}
		}
		c.String(http.StatusInternalServerError, "")
	}
}
func GetTeachCourse(c *gin.Context) {
	if claims, ok := c.Get("claims"); ok {
		s := service.CourseService{
			UserId: claims.(*utils.Claims).Id,
		}
		if courses, err := s.GetCourseByUserId(); err == nil {
			c.JSON(http.StatusOK, courses)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
