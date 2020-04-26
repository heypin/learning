package api

import (
	"github.com/gin-gonic/gin"
	"learning/models"
	"learning/service"
	"log"
	"net/http"
	"strconv"
)

type CreateExamLibItemForm struct {
	ExamLibId uint     `json:"examLibId" binding:"required"`
	Type      string   `json:"type" binding:"required"`
	Question  string   `json:"question" binding:"required"`
	Answer    string   `json:"answer" binding:"required"`
	Score     uint     `json:"score" binding:"required"`
	Options   []Option `json:"options" `
}

func CreateExamLibItemAndOptions(c *gin.Context) {
	var form CreateExamLibItemForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	options := make([]*models.ExamLibItemOption, 0)
	if form.Type == models.Subject_Single || form.Type == models.Subject_Multiple {
		for _, v := range form.Options {
			options = append(options, &models.ExamLibItemOption{
				Sequence: v.Sequence,
				Content:  v.Content,
			})
		}
	}
	s := service.ExamLibItemService{
		ExamLibId: form.ExamLibId,
		Type:      form.Type,
		Question:  form.Question,
		Answer:    form.Answer,
		Score:     form.Score,
		Options:   options,
	}
	if _, err := s.CreateExamLibItemAndOptions(); err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusCreated, "")
	}
}

type UpdateExamLibItemForm struct {
	Id        uint     `json:"id" binding:"required"`
	ExamLibId uint     `json:"examLibId" binding:"required"`
	Type      string   `json:"type" binding:"required"`
	Question  string   `json:"question" binding:"required"`
	Answer    string   `json:"answer" binding:"required"`
	Score     uint     `json:"score" binding:"required"`
	Options   []Option `json:"options" `
}

func UpdateExamLibItemAndOptions(c *gin.Context) {
	var form UpdateExamLibItemForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	options := make([]*models.ExamLibItemOption, 0)
	if form.Type == models.Subject_Single || form.Type == models.Subject_Multiple {
		for _, v := range form.Options {
			options = append(options, &models.ExamLibItemOption{
				Sequence: v.Sequence,
				Content:  v.Content,
			})
		}
	}
	s := service.ExamLibItemService{
		Id:        form.Id,
		ExamLibId: form.ExamLibId,
		Type:      form.Type,
		Question:  form.Question,
		Answer:    form.Answer,
		Score:     form.Score,
		Options:   options,
	}
	if err := s.UpdateExamLibItemAndOptions(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetExamLibItemsByLibId(c *gin.Context) {
	examLibId, err := strconv.Atoi(c.Query("examLibId"))
	if err != nil || examLibId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibItemService{
		ExamLibId: uint(examLibId),
	}
	if items, err := s.GetExamLibItemsByLibId(); err == nil {
		c.JSON(http.StatusOK, items)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func DeleteExamLibItemById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamLibItemService{
		Id: uint(id),
	}
	if err := s.DeleteExamLibItemAndOptions(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
