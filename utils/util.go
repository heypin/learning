package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mapset "github.com/deckarep/golang-set"
	"gopkg.in/gomail.v2"
	"io"
	"learning/conf"
	"learning/models"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateClassCode(id uint) string { //生成班级码
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1<<16 - 1)
	return fmt.Sprintf("%04X%02X", random, id)
}

//发送邮件
func SendEmail(to string, content string) error {
	email := conf.AppConfig.Email
	m := gomail.NewMessage()
	m.SetHeader("Subject", "[辅助学习平台]")
	m.SetHeader("From", email.Username)
	m.SetHeader("To", to)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(email.Host, email.Port,
		email.Username, email.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
func SetMarkAndScore(subjectType string, rightAnswer string, totalScore uint, userAnswer string, userScore *uint, mark *uint) {
	if subjectType == models.SubjectShort || subjectType == models.SubjectProgram { //如果有主观题标为未评
		*mark = 0
	} else if subjectType == models.SubjectSingle || subjectType == models.SubjectJudgement {
		if userAnswer == rightAnswer {
			*userScore = totalScore
		}
	} else if subjectType == models.SubjectMultiple {
		submitSet := mapset.NewSet()
		for _, v := range strings.Split(userAnswer, ",") {
			submitSet.Add(v)
		}
		rightSet := mapset.NewSet()
		for _, v := range strings.Split(rightAnswer, ",") {
			rightSet.Add(v)
		}
		if submitSet.Equal(rightSet) {
			*userScore = totalScore
		}
	} else if subjectType == models.SubjectBlank {
		rightArr := strings.Split(rightAnswer, ",")
		submitArr := strings.Split(userAnswer, ",")
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
		*userScore = uint(float32(rightCount) / float32(len(rightArr)) * float32(totalScore))
	}
}

type Option struct {
	Sequence string
	Content  string
}
type Subject struct {
	Type     string
	Question string
	Answer   string
	Score    uint
	Options  []Option
}

func isSupportSubjectType(subjectType string) bool {
	set := mapset.NewSet()
	set.Add(models.SubjectSingle)
	set.Add(models.SubjectMultiple)
	set.Add(models.SubjectJudgement)
	set.Add(models.SubjectBlank)
	set.Add(models.SubjectShort)
	set.Add(models.SubjectProgram)
	return set.Contains(subjectType)
}
func ReadExcelToSubject(reader io.Reader) ([]Subject, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	subjects := make([]Subject, 0)
	for index, row := range rows {
		if index == 0 { //跳过题目行
			continue
		}
		var item Subject
		subjectType := strings.TrimSpace(row[0])
		if !isSupportSubjectType(subjectType) { //跳过不支持的题目类型
			continue
		}
		item.Type = subjectType
		item.Question = row[1]
		item.Answer = strings.TrimSpace(row[2])
		if score, err := strconv.Atoi(strings.TrimSpace(row[3])); err != nil {
			continue
		} else {
			item.Score = uint(score)
		}
		if subjectType == models.SubjectSingle ||
			subjectType == models.SubjectMultiple {
			item.Answer = strings.ToUpper(item.Answer)
			item.Options = make([]Option, 0)
			for i := 4; i < len(row); i++ {
				item.Options = append(item.Options, Option{
					Sequence: string([]byte{byte(65 - 4 + i)}), //65转为字符串为A,66为B
					Content:  row[i],
				})
			}
		}
		subjects = append(subjects, item)
	}
	return subjects, nil
}
