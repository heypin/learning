package api

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"learning/cache"
	"learning/conf"
	"learning/service"
	"learning/utils"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func PlayVideo(c *gin.Context) {
	videoName := c.Param("name")
	localPath := conf.AppConfig.Path.Video + "/" + videoName
	video, err := os.Open(localPath)
	defer func() {
		if err := video.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.Header("Content-Type", "video/mp4")
	http.ServeContent(c.Writer, c.Request, "", time.Now(), video)
}

type ExecuteProgramForm struct {
	Language string `json:"language" binding:"required"`
	Input    string `json:"input" binding:"required"`
}

func ExecuteProgram(c *gin.Context) {
	var form ExecuteProgramForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	if out, err := utils.ExecuteProgramSubject(form.Language, form.Input); err == nil {
		c.JSON(http.StatusOK, out)
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}
}
func ExportExamToExcel(c *gin.Context) {
	examPublishId, err := strconv.Atoi(c.Query("examPublishId"))
	if err != nil || examPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ExamSubmitService{
		ExamPublishId: uint(examPublishId),
	}
	submits, err := s.GetExamSubmitsByPublishId()
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	title := map[string]string{
		"A1": "帐号", "B1": "姓名", "C1": "学号", "D1": "总分", "E1": "开考时间", "F1": "完成时间",
	}
	for k, v := range title {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	content := make(map[string]interface{})
	for i, v := range submits {
		content["A"+strconv.Itoa(i+2)] = v.User.Email
		content["B"+strconv.Itoa(i+2)] = v.User.RealName
		content["C"+strconv.Itoa(i+2)] = v.User.Number
		content["D"+strconv.Itoa(i+2)] = v.TotalScore
		content["E"+strconv.Itoa(i+2)] = strings.Split(v.StartTime.String(), "+")[0]
		if v.FinishTime == nil {
			content["F"+strconv.Itoa(i+2)] = "截止时间"
		} else {
			content["F"+strconv.Itoa(i+2)] = strings.Split(v.FinishTime.String(), "+")[0]
		}
	}
	for k, v := range content {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	ps := service.ExamPublishService{
		Id: uint(examPublishId),
	}
	var filename string
	if publish, _ := ps.GetExamPublishById(); publish != nil {
		filename = publish.ExamLib.Name
	}
	filename = url.QueryEscape(filename) //防止中文乱码
	_ = f.SetColWidth("Sheet1", "A", "F", 20)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	if err = f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
}

func ExportHomeworkToExcel(c *gin.Context) {
	homeworkPublishId, err := strconv.Atoi(c.Query("homeworkPublishId"))
	if err != nil || homeworkPublishId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.HomeworkSubmitService{
		HomeworkPublishId: uint(homeworkPublishId),
	}
	submits, err := s.GetHomeworkSubmitsByPublishId()
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	title := map[string]string{
		"A1": "帐号", "B1": "姓名", "C1": "学号", "D1": "总分", "E1": "首次提交时间", "F1": "更新时间",
	}
	for k, v := range title {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	content := make(map[string]interface{})
	for i, v := range submits {
		content["A"+strconv.Itoa(i+2)] = v.User.Email
		content["B"+strconv.Itoa(i+2)] = v.User.RealName
		content["C"+strconv.Itoa(i+2)] = v.User.Number
		content["D"+strconv.Itoa(i+2)] = v.TotalScore
		content["E"+strconv.Itoa(i+2)] = strings.Split(v.CreatedAt.String(), "+")[0]
		content["F"+strconv.Itoa(i+2)] = strings.Split(v.UpdatedAt.String(), "+")[0]
	}
	for k, v := range content {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	ps := service.HomeworkPublishService{
		Id: uint(homeworkPublishId),
	}
	var filename string
	if publish, _ := ps.GetHomeworkPublishById(); publish != nil {
		filename = publish.HomeworkLib.Name
	}
	filename = url.QueryEscape(filename) //防止中文乱码
	_ = f.SetColWidth("Sheet1", "A", "F", 20)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	if err = f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
}
func GenerateRegisterCode(c *gin.Context) {
	form := struct {
		Email string `form:"email" binding:"required,email"`
	}{}
	if err := c.ShouldBindQuery(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(9999)
	registerCode := fmt.Sprintf("%04d", random)
	err := cache.RedisClient.Set(cache.CaptchaPrefix+"."+form.Email, registerCode, time.Minute*5).Err()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	if err := utils.SendEmail(form.Email, registerCode); err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
}
