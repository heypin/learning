package api

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"learning/cache"
	"learning/conf"
	"learning/models"
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
func DownloadExcelExample(c *gin.Context) {
	filename := url.QueryEscape("example.xlsx")
	filepath := "./static/example/example.xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(filepath)
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
func ExportExamSubmitToExcel(c *gin.Context) {
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

func ExportHomeworkSubmitToExcel(c *gin.Context) {
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
func GenerateCaptcha(c *gin.Context) {
	form := struct {
		Email string `form:"email" binding:"required,email"`
	}{}
	if err := c.ShouldBindQuery(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "请输入邮箱"})
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
	content := fmt.Sprintf("你的注册验证码为<b>%s</b>，五分钟内有效", registerCode)
	if err := utils.SendEmail(form.Email, content); err != nil {
		c.JSON(http.StatusInternalServerError, "邮件发送失败")
		return
	}
}

type ImportAndExportSubjectForm struct {
	LibId uint   `form:"libId" binding:"required"`
	Type  string `form:"type" binding:"required"` //是试题库还是作业库
}

func ExportLibSubjectToExcel(c *gin.Context) {
	var form ImportAndExportSubjectForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	title := map[string]string{
		"A1": "题目类型", "B1": "问题", "C1": "答案", "D1": "分数", "E1": "选项A", "F1": "选项B...",
	}
	for k, v := range title {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	content := make(map[string]interface{})
	var filename = ""
	if form.Type == "exam" {
		s := service.ExamLibService{Id: form.LibId}
		examLib, err := s.GetExamLibWithItemsById()
		if err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}
		filename = examLib.Name
		for i, v := range examLib.Items {
			content["A"+strconv.Itoa(i+2)] = v.Type
			content["B"+strconv.Itoa(i+2)] = v.Question
			content["C"+strconv.Itoa(i+2)] = *v.Answer
			content["D"+strconv.Itoa(i+2)] = v.Score
			if v.Type == models.SubjectSingle ||
				v.Type == models.SubjectMultiple {
				for j, option := range v.Options {
					sequence := string([]byte{byte(69 + j)}) //69转为字符串为E,70为F
					content[sequence+strconv.Itoa(i+2)] = option.Content
				}
			}
		}
	} else if form.Type == "homework" {
		s := service.HomeworkLibService{Id: form.LibId}
		homeworkLib, err := s.GetHomeworkLibWithItemsById()
		if err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}
		filename = homeworkLib.Name
		for i, v := range homeworkLib.Items {
			content["A"+strconv.Itoa(i+2)] = v.Type
			content["B"+strconv.Itoa(i+2)] = v.Question
			content["C"+strconv.Itoa(i+2)] = *v.Answer
			content["D"+strconv.Itoa(i+2)] = v.Score
			if v.Type == models.SubjectSingle ||
				v.Type == models.SubjectMultiple {
				for j, option := range v.Options {
					sequence := string([]byte{byte(69 + j)}) //69转为字符串为E,70为F
					content[sequence+strconv.Itoa(i+2)] = option.Content
				}
			}
		}
	}

	for k, v := range content {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	filename = url.QueryEscape(filename) //防止中文乱码
	_ = f.SetColWidth("Sheet1", "A", "D", 20)
	_ = f.SetColWidth("Sheet1", "E", "H", 40)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
}

func ImportExcelSubjectToLib(c *gin.Context) {
	var form ImportAndExportSubjectForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if fileHeader, err := c.FormFile("excel"); err == nil {
		if file, err := fileHeader.Open(); err == nil {
			if subjects, err := utils.ReadExcelToSubject(file); err == nil {
				for _, subject := range subjects {
					if form.Type == "exam" {
						itemService := service.ExamLibItemService{
							ExamLibId: form.LibId,
							Type:      subject.Type,
							Question:  subject.Question,
							Answer:    &subject.Answer,
							Score:     subject.Score,
						}
						if subject.Options != nil {
							for _, option := range subject.Options {
								itemService.Options = append(itemService.Options, &models.ExamLibItemOption{
									Sequence: option.Sequence,
									Content:  option.Content,
								})
							}
						}
						_, _ = itemService.CreateExamLibItemAndOptions()
					} else if form.Type == "homework" {
						itemService := service.HomeworkLibItemService{
							HomeworkLibId: form.LibId,
							Type:          subject.Type,
							Question:      subject.Question,
							Answer:        &subject.Answer,
							Score:         subject.Score,
						}
						if subject.Options != nil {
							for _, option := range subject.Options {
								itemService.Options = append(itemService.Options, &models.HomeworkLibItemOption{
									Sequence: option.Sequence,
									Content:  option.Content,
								})
							}
						}
						_, _ = itemService.CreateLibItemAndOptions()
					}
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"err": "读取题目失败",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "文件打开失败",
			})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "文件获取失败",
		})
	}
}
