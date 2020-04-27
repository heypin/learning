package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/service"
	"learning/utils"
	"net/http"
	"strconv"
	"time"
)

func StartExam(c *gin.Context) {
	form := struct {
		ExamPublishId uint `json:"examPublishId" binding:"required"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.ExamSubmitService{
			UserId:        claims.(*utils.Claims).Id,
			ExamPublishId: form.ExamPublishId,
		}
		if submit, err := s.CreateExamSubmit(); err == nil {
			c.JSON(http.StatusOK, submit)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

func FinishExam(c *gin.Context) {
	form := struct {
		Id uint `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	now := time.Now()
	s := service.ExamSubmitService{
		Id:         form.Id,
		FinishTime: &now,
	}
	if err := s.UpdateExamSubmitById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamSubmitById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamSubmitService{
		Id: uint(id),
	}
	if submit, err := s.GetExamSubmitById(); err == nil {
		c.JSON(http.StatusOK, submit)
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}
}
func GetExamSubmitsByPublishId(c *gin.Context) {
	examPublishId, err := strconv.Atoi(c.Query("examPublishId"))
	if err != nil || examPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamSubmitService{
		ExamPublishId: uint(examPublishId),
	}
	if submits, err := s.GetExamSubmitsByPublishId(); err == nil {
		c.JSON(http.StatusOK, submits)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamUserSubmitWithItems(c *gin.Context) {
	examPublishId, err := strconv.Atoi(c.Query("examPublishId"))
	if err != nil || examPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamSubmitService{
		ExamPublishId: uint(examPublishId),
	}
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil || userId <= 0 {
		if claims, ok := c.Get("claims"); ok {
			s.UserId = claims.(*utils.Claims).Id
		} else {
			c.String(http.StatusInternalServerError, "")
			return
		}
	} else {
		s.UserId = uint(userId)
	}
	if submit, err := s.GetExamUserSubmitWithItems(); err == nil {
		c.JSON(http.StatusOK, submit)
		return
	}
	c.String(http.StatusInternalServerError, "")
}

type UpdateExamSubmitItem struct {
	Id            uint  `json:"id"`
	ExamLibItemId uint  `json:"examLibItemId" binding:"required"`
	Score         *uint `json:"score"`
}
type UpdateExamSubmitForm struct {
	Id          uint                   `json:"id" binding:"required"`
	SubmitItems []UpdateExamSubmitItem `json:"submitItems"`
}

func UpdateExamSubmitItemsScore(c *gin.Context) {
	var form UpdateExamSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	submitItems := make([]*models.ExamSubmitItem, 0)
	for _, item := range form.SubmitItems {
		submitItems = append(submitItems, &models.ExamSubmitItem{
			Model:         gorm.Model{ID: item.Id},
			ExamSubmitId:  form.Id,
			ExamLibItemId: item.ExamLibItemId,
			Score:         item.Score,
		})
	}
	var mark uint = 1
	s := service.ExamSubmitService{
		Id:          form.Id,
		Mark:        &mark,
		SubmitItems: submitItems,
	}
	if err := s.UpdateExamSubmitWithItems(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type ExamSubmitItem struct {
	Id            uint   `json:"id"`
	ExamLibItemId uint   `json:"examLibItemId" binding:"required"`
	Answer        string `json:"answer"`
}
type ExamSubmitForm struct {
	Id            uint             `json:"id" binding:"required"`
	ExamPublishId uint             `json:"examPublishId" binding:"required"`
	SubmitItems   []ExamSubmitItem `json:"submitItems"`
}

func SubmitExamItem(c *gin.Context) {
	var form ExamSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	publishService := service.ExamPublishService{
		Id: form.ExamPublishId,
	}
	claims, _ := c.Get("claims")
	submitService := service.ExamSubmitService{
		Id:            form.Id,
		ExamPublishId: form.ExamPublishId,
		UserId:        claims.(*utils.Claims).Id,
	}
	if publish, err := publishService.GetExamPublishById(); err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	} else {
		now := time.Now()
		if now.Before(publish.BeginTime) || now.After(publish.EndTime) {
			c.String(http.StatusBadRequest, "")
			return
		}
		submit, _ := submitService.GetUserExamSubmitRecord()
		if submit == nil {
			if publish.EndTime.Before(now) {
				c.String(http.StatusBadRequest, "")
				return
			}
		} else {
			if now.Sub(submit.StartTime).Minutes() > float64(publish.Duration) ||
				submit.FinishTime != nil {
				c.String(http.StatusBadRequest, "")
				return
			}
		}
	}
	submitItems := make([]*models.ExamSubmitItem, 0)
	for _, item := range form.SubmitItems {
		submitItems = append(submitItems, &models.ExamSubmitItem{
			Model:         gorm.Model{ID: item.Id},
			ExamLibItemId: item.ExamLibItemId,
			Answer:        item.Answer,
			Score:         new(uint),
		})
	}
	var mark uint = 1
	submitService.Mark = &mark
	submitService.SubmitItems = submitItems
	for _, submitItem := range submitService.SubmitItems {
		s := service.ExamLibItemService{
			Id: submitItem.ExamLibItemId,
		}
		if libItem, err := s.GetExamLibItemById(); err == nil && libItem != nil {
			utils.SetMarkAndScore(libItem.Type, libItem.Answer, libItem.Score,
				submitItem.Answer, submitItem.Score, &mark)
		}
	}
	if err := submitService.UpdateExamSubmitWithItems(); err == nil {
		c.String(http.StatusOK, "")
		return
	}
	c.String(http.StatusInternalServerError, "")
}
