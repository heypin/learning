package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/service"
	"learning/utils"
	"net/http"
	"strconv"
	"strings"
)

func GetHomeworkSubmitsByPublishId(c *gin.Context) {
	homeworkPublishId, err := strconv.Atoi(c.Query("homeworkPublishId"))
	if err != nil || homeworkPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkSubmitService{
		HomeworkPublishId: uint(homeworkPublishId),
	}
	if submits, err := s.GetHomeworkSubmitsByPublishId(); err == nil {
		c.JSON(http.StatusOK, submits)
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}

}
func GetHomeworkSubmitWithItems(c *gin.Context) {
	homeworkPublishId, err := strconv.Atoi(c.Query("homeworkPublishId"))
	if err != nil || homeworkPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.HomeworkSubmitService{
			UserId:            claims.(*utils.Claims).Id,
			HomeworkPublishId: uint(homeworkPublishId),
		}
		if submit, err := s.GetHomeworkSubmitWithItems(); err == nil {
			c.JSON(http.StatusOK, submit)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

type SubmitItem struct {
	Id                uint   `form:"id"`
	HomeworkLibItemId uint   `form:"homeworkLibItemId" binding:"required"`
	Answer            string `form:"answer"`
}
type HomeworkSubmitForm struct {
	Id                uint         `form:"id"`
	HomeworkPublishId uint         `form:"homeworkPublishId" binding:"required"`
	SubmitItems       []SubmitItem `form:"submitItems"`
}

func SubmitHomeworkWithItems(c *gin.Context) {
	var form HomeworkSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	submitItems := make([]*models.HomeworkSubmitItem, 0)
	for _, item := range form.SubmitItems {
		submitItems = append(submitItems, &models.HomeworkSubmitItem{
			Model:             gorm.Model{ID: item.Id},
			HomeworkLibItemId: item.HomeworkLibItemId,
			Answer:            item.Answer,
		})
	}
	if claims, ok := c.Get("claims"); ok {
		var mark uint = 1
		submitService := service.HomeworkSubmitService{
			UserId:            claims.(*utils.Claims).Id,
			HomeworkPublishId: form.HomeworkPublishId,
			Mark:              &mark,
			SubmitItems:       submitItems,
		}
		for _, submitItem := range submitService.SubmitItems {
			s := service.HomeworkLibItemService{
				Id: submitItem.HomeworkLibItemId,
			}
			if libItem, err := s.GetHomeworkLibItemById(); err != nil && libItem != nil {
				setScore(submitItem, libItem, &mark)
				submitService.TotalScore += submitItem.Score
			}
		}
		if _, err := submitService.SubmitHomeworkWithItems(); err == nil {
			c.String(http.StatusOK, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
func setScore(submitItem *models.HomeworkSubmitItem, libItem *models.HomeworkLibItem, mark *uint) {
	if libItem.Type == models.Subject_Short || libItem.Type == models.Subject_Program { //如果有主观题标为未评
		*mark = 0
	} else if libItem.Type == models.Subject_Single ||
		libItem.Type == models.Subject_Multiple ||
		libItem.Type == models.Subject_Judgement {

		if submitItem.Answer == libItem.Answer {
			submitItem.Score = libItem.Score
		}
	} else if libItem.Type == models.Subject_Blank {
		rightArr := strings.Split(libItem.Answer, ",")
		submitArr := strings.Split(submitItem.Answer, ",")
		var length int
		if len(submitArr) < len(rightArr) {
			length = len(submitArr)
		} else {
			length = len(rightArr)
		}
		var rightCount int
		for i := 0; i < length; i++ {
			if strings.TrimSpace(submitArr[i]) == strings.TrimSpace(rightArr[i]) {
				rightCount++
			}
		}
		submitItem.Score = uint(rightCount / len(rightArr))
	}
}
