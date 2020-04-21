package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PublishHomeworkForm struct {
	ClassId       uint      `form:"classId" binding:"required"`
	HomeworkLibId uint      `form:"homeworkLibId" binding:"required"`
	BeginTime     time.Time `form:"beginTime" binding:"required"`
	EndTime       time.Time `form:"endTime" binding:"required"`
	Resubmit      *uint     `form:"resubmit" binding:"required"`
}

func PublishHomework(c *gin.Context) {
	var form PublishHomeworkForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "")
		return
	}
	if form.BeginTime.After(form.EndTime) {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkPublishService{
		ClassId:       form.ClassId,
		HomeworkLibId: form.HomeworkLibId,
		BeginTime:     form.BeginTime,
		EndTime:       form.EndTime,
		Resubmit:      form.Resubmit,
	}

	if id, err := s.PublishHomework(); err == nil && id != 0 {
		c.String(http.StatusCreated, "")
	} else if err == nil && id == 0 {
		c.JSON(http.StatusAccepted, gin.H{"code": http.StatusAccepted})
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type UpdateHomeworkPublishForm struct {
	Id        uint      `form:"id" binding:"required"`
	BeginTime time.Time `form:"beginTime" binding:"required"`
	EndTime   time.Time `form:"endTime" binding:"required"`
	Resubmit  *uint     `form:"resubmit" binding:"required"`
}

func UpdateHomeworkPublishById(c *gin.Context) {
	var form UpdateHomeworkPublishForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkPublishService{
		Id:        form.Id,
		BeginTime: form.BeginTime,
		EndTime:   form.EndTime,
		Resubmit:  form.Resubmit,
	}
	if err := s.UpdateHomeworkPublishById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetHomeworkPublishesByClassId(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkPublishService{
		ClassId: uint(classId),
	}
	if publishes, err := s.GetHomeworkPublishesByClassId(); err == nil {
		c.JSON(http.StatusOK, publishes)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetHomeworkPublishesWithSubmitByClassId(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.HomeworkPublishService{
			ClassId: uint(classId),
		}
		userId := claims.(*utils.Claims).Id
		if publishes, err := s.GetHomeworkPublishesWithSubmitByClassId(userId); err == nil {
			c.JSON(http.StatusOK, publishes)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
