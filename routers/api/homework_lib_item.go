package api

import (
	"github.com/gin-gonic/gin"
	"learning/models"
	"learning/service"
	"log"
	"net/http"
	"strconv"
)

type Option struct {
	Sequence string `json:"sequence" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
type CreateLibItemForm struct {
	HomeworkLibId uint     `json:"homeworkLibId" binding:"required"`
	Type          string   `json:"type" binding:"required"`
	Question      string   `json:"question" binding:"required"`
	Answer        string   `json:"answer" binding:"required"`
	Score         uint     `json:"score" binding:"required"`
	Options       []Option `json:"options" `
}

func CreateHomeworkLibItemAndOptions(c *gin.Context) {
	var form CreateLibItemForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	var options []*models.HomeworkLibItemOption = nil
	if form.Type == models.Subject_Single || form.Type == models.Subject_Multiple {
		options = []*models.HomeworkLibItemOption{}
		for _, v := range form.Options {
			options = append(options, &models.HomeworkLibItemOption{
				Sequence: v.Sequence,
				Content:  v.Content,
			})
		}
	}
	s := service.HomeworkLibItemService{
		HomeworkLibId: form.HomeworkLibId,
		Type:          form.Type,
		Question:      form.Question,
		Answer:        form.Answer,
		Score:         form.Score,
		Options:       options,
	}
	if _, err := s.CreateLibItemAndOptions(); err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusCreated, "")
	}
}

type UpdateLibItemForm struct {
	Id            uint     `json:"id" binding:"required"`
	HomeworkLibId uint     `json:"homeworkLibId" binding:"required"`
	Type          string   `json:"type" binding:"required"`
	Question      string   `json:"question" binding:"required"`
	Answer        string   `json:"answer" binding:"required"`
	Score         uint     `json:"score" binding:"required"`
	Options       []Option `json:"options" `
}

func UpdateHomeworkLibItemAndOptions(c *gin.Context) {
	var form UpdateLibItemForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	var options []*models.HomeworkLibItemOption = nil
	if form.Type == models.Subject_Single || form.Type == models.Subject_Multiple {
		options = []*models.HomeworkLibItemOption{}
		for _, v := range form.Options {
			options = append(options, &models.HomeworkLibItemOption{
				Sequence: v.Sequence,
				Content:  v.Content,
			})
		}
	}
	s := service.HomeworkLibItemService{
		Id:            form.Id,
		HomeworkLibId: form.HomeworkLibId,
		Type:          form.Type,
		Question:      form.Question,
		Answer:        form.Answer,
		Score:         form.Score,
		Options:       options,
	}
	if err := s.UpdateLibItemAndOptions(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetHomeworkLibItemsByLibId(c *gin.Context) {
	homeworkLibId, err := strconv.Atoi(c.Query("homeworkLibId"))
	if err != nil || homeworkLibId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibItemService{
		HomeworkLibId: uint(homeworkLibId),
	}
	if items, err := s.GetLibItemsByLibId(); err == nil {
		c.JSON(http.StatusOK, items)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func DeleteHomeworkLibItemById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkLibItemService{
		Id: uint(id),
	}
	if err := s.DeleteLibItemAndOptionsById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
