package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"net/http"
	"strconv"
	"time"
)

type PublishExamForm struct {
	ClassId   uint      `json:"classId" binding:"required"`
	ExamLibId uint      `json:"examLibId" binding:"required"`
	BeginTime time.Time `json:"beginTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Duration  uint      `json:"duration" binding:"required"`
}

func PublishExam(c *gin.Context) {
	var form PublishExamForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if form.EndTime.Sub(form.BeginTime).Minutes() < float64(form.Duration) {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamPublishService{
		ClassId:   form.ClassId,
		ExamLibId: form.ExamLibId,
		BeginTime: form.BeginTime,
		EndTime:   form.EndTime,
		Duration:  form.Duration,
	}
	if id, err := s.PublishExam(); err == nil && id != 0 {
		c.String(http.StatusCreated, "")
	} else if err == nil && id == 0 {
		c.JSON(http.StatusAccepted, gin.H{"code": http.StatusAccepted})
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type UpdateExamPublishForm struct {
	Id        uint      `json:"id" binding:"required"`
	BeginTime time.Time `json:"beginTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Duration  uint      `json:"duration" binding:"required"`
}

func UpdateExamPublishById(c *gin.Context) {
	var form UpdateExamPublishForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamPublishService{
		Id:        form.Id,
		BeginTime: form.BeginTime,
		EndTime:   form.EndTime,
		Duration:  form.Duration,
	}
	if err := s.UpdateExamPublishById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamPublishById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamPublishService{
		Id: uint(id),
	}
	if publish, err := s.GetExamPublishById(); err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, publish)
	}
}
func GetExamPublishesByClassId(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamPublishService{
		ClassId: uint(classId),
	}
	if publishes, err := s.GetExamPublishesByClassId(); err == nil {
		c.JSON(http.StatusOK, publishes)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamPublishesWithSubmitByClassId(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.ExamPublishService{
			ClassId: uint(classId),
		}
		userId := claims.(*utils.Claims).Id
		if publishes, err := s.GetExamPublishesWithUserSubmitByClassId(userId); err == nil {
			c.JSON(http.StatusOK, publishes)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
