package api

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetHomeworkSubmitById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkSubmitService{
		Id: uint(id),
	}
	if submit, err := s.GetHomeworkSubmitById(); err == nil {
		c.JSON(http.StatusOK, submit)
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}
}
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

type UpdateHomeworkMarkForm struct {
	Id   uint  `json:"id" binding:"required"`
	Mark *uint `json:"mark" binding:"required"`
}

func UpdateHomeworkSubmitMarkById(c *gin.Context) {
	var form UpdateHomeworkMarkForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkSubmitService{
		Id:   form.Id,
		Mark: form.Mark,
	}
	if err := s.UpdateHomeworkSubmitMarkById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetHomeworkUserSubmitWithItems(c *gin.Context) {
	homeworkPublishId, err := strconv.Atoi(c.Query("homeworkPublishId"))
	if err != nil || homeworkPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkSubmitService{
		HomeworkPublishId: uint(homeworkPublishId),
	}
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil || userId <= 0 {
		s.UserId = uint(userId)
	} else {
		if claims, ok := c.Get("claims"); ok {
			s.UserId = claims.(*utils.Claims).Id
		} else {
			c.String(http.StatusInternalServerError, "")
			return
		}
	}
	if submit, err := s.GetHomeworkUserSubmitWithItems(); err == nil {
		c.JSON(http.StatusOK, submit)
		return
	}
	c.String(http.StatusInternalServerError, "")
}

type UpdateSubmitItem struct {
	Id                uint `form:"id"`
	HomeworkLibItemId uint `form:"homeworkLibItemId" binding:"required"`
	Score             uint `form:"score"`
}
type UpdateHomeworkSubmitForm struct {
	Id          uint               `form:"id" binding:"required"`
	SubmitItems []UpdateSubmitItem `form:"submitItems"`
}

func UpdateHomeworkSubmitItemsScore(c *gin.Context) {
	var form UpdateHomeworkSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	submitItems := make([]*models.HomeworkSubmitItem, 0)
	for _, item := range form.SubmitItems {
		submitItems = append(submitItems, &models.HomeworkSubmitItem{
			Model:             gorm.Model{ID: item.Id},
			HomeworkSubmitId:  form.Id,
			HomeworkLibItemId: item.HomeworkLibItemId,
			Score:             item.Score,
		})
	}
	var mark uint = 1
	s := service.HomeworkSubmitService{
		Id:          form.Id,
		Mark:        &mark,
		SubmitItems: submitItems,
	}
	if err := s.UpdateSubmitHomeworkWithItems(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

type SubmitItem struct {
	Id                uint   `form:"id"`
	HomeworkLibItemId uint   `form:"homeworkLibItemId" binding:"required"`
	Answer            string `form:"answer"`
}
type HomeworkSubmitForm struct {
	Id                uint         `form:"id" ` //
	HomeworkPublishId uint         `form:"homeworkPublishId" binding:"required"`
	SubmitItems       []SubmitItem `form:"submitItems"`
}

func SubmitHomeworkWithItems(c *gin.Context) {
	var form HomeworkSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		log.Println(err)
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
			Id:                form.Id,
			UserId:            claims.(*utils.Claims).Id,
			HomeworkPublishId: form.HomeworkPublishId,
			Mark:              &mark,
			SubmitItems:       submitItems,
		}
		for _, submitItem := range submitService.SubmitItems {
			s := service.HomeworkLibItemService{
				Id: submitItem.HomeworkLibItemId,
			}
			if libItem, err := s.GetHomeworkLibItemById(); err == nil && libItem != nil {
				setScore(submitItem, libItem, &mark)
				//submitService.TotalScore += submitItem.Score 不在这里计算，在数据库中计算出总分
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
		libItem.Type == models.Subject_Judgement {
		if submitItem.Answer == libItem.Answer {
			submitItem.Score = libItem.Score
		}
	} else if libItem.Type == models.Subject_Multiple {
		submitSet := mapset.NewSet()
		for _, v := range strings.Split(submitItem.Answer, ",") {
			submitSet.Add(v)
		}
		rightSet := mapset.NewSet()
		for _, v := range strings.Split(libItem.Answer, ",") {
			rightSet.Add(v)
		}
		if submitSet.Equal(rightSet) {
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
