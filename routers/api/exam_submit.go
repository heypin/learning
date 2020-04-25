package api

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/service"
	"learning/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StartExamForm struct {
	ExamPublishId uint `form:"examPublishId" binding:"required"`
}

func StartExam(c *gin.Context) {
	var form StartExamForm
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

type EndExamForm struct {
	Id            uint `form:"id" binding:"required"`
	ExamPublishId uint `form:"examPublishId" binding:"required"`
}

func FinishExam(c *gin.Context) {
	var form EndExamForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamSubmitService{
		Id:         form.Id,
		FinishTime: time.Now(),
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
		s.UserId = uint(userId)
	} else {
		if claims, ok := c.Get("claims"); ok {
			s.UserId = claims.(*utils.Claims).Id
		} else {
			c.String(http.StatusInternalServerError, "")
			return
		}
	}
	if submit, err := s.GetExamUserSubmitWithItems(); err == nil {
		c.JSON(http.StatusOK, submit)
		return
	}
	c.String(http.StatusInternalServerError, "")
}

type UpdateExamSubmitItem struct {
	Id            uint `form:"id"`
	ExamLibItemId uint `form:"examLibItemId" binding:"required"`
	Score         uint `form:"score"`
}
type UpdateExamSubmitForm struct {
	Id          uint                   `form:"id" binding:"required"`
	SubmitItems []UpdateExamSubmitItem `form:"submitItems"`
}

func UpdateExamSubmitItemsScore(c *gin.Context) {
	var form UpdateExamSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	submitItems := make([]*models.ExamSubmitItem, 5)
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
	Id            uint   `form:"id"`
	ExamLibItemId uint   `form:"examLibItemId" binding:"required"`
	Answer        string `form:"answer"`
}
type ExamSubmitForm struct {
	Id            uint             `form:"id" binding:"required"` //
	ExamPublishId uint             `form:"examPublishId" binding:"required"`
	SubmitItems   []ExamSubmitItem `form:"submitItems"`
}

func SubmitExamItem(c *gin.Context) {
	var form ExamSubmitForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	submitItems := make([]*models.ExamSubmitItem, 1)
	for _, item := range form.SubmitItems {
		submitItems = append(submitItems, &models.ExamSubmitItem{
			Model:         gorm.Model{ID: item.Id},
			ExamLibItemId: item.ExamLibItemId,
			Answer:        item.Answer,
		})
	}
	var mark uint = 1
	submitService := service.ExamSubmitService{
		Id:          form.Id,
		Mark:        &mark,
		SubmitItems: submitItems,
	}
	for _, submitItem := range submitService.SubmitItems {
		s := service.ExamLibItemService{
			Id: submitItem.ExamLibItemId,
		}
		if libItem, err := s.GetExamLibItemById(); err == nil && libItem != nil {
			setExamScore(submitItem, libItem, &mark)
		}
	}
	if err := submitService.UpdateExamSubmitWithItems(); err == nil {
		c.String(http.StatusOK, "")
		return
	}
	c.String(http.StatusInternalServerError, "")
}
func setExamScore(submitItem *models.ExamSubmitItem, libItem *models.ExamLibItem, mark *uint) {
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
