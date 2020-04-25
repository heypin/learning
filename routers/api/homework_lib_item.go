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
	Sequence string `form:"sequence" binding:"required"`
	Content  string `form:"content" binding:"required"`
}
type CreateLibItemForm struct {
	HomeworkLibId uint     `form:"homeworkLibItemId" binding:"required"`
	Type          string   `form:"type" binding:"required"`
	Question      string   `form:"question" binding:"required"`
	Answer        string   `form:"answer" binding:"required"`
	Score         uint     `form:"score" binding:"required"`
	Options       []Option `form:"options" `
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
	Id       uint     `form:"id" binding:"required"`
	Type     string   `form:"type" binding:"required"`
	Question string   `form:"question" binding:"required"`
	Answer   string   `form:"answer" binding:"required"`
	Score    uint     `form:"score" binding:"required"`
	Options  []Option `form:"options" `
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
		Id:       form.Id,
		Type:     form.Type,
		Question: form.Question,
		Answer:   form.Answer,
		Score:    form.Score,
		Options:  options,
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
