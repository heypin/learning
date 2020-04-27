package utils

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"learning/models"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GenerateClassCode(id uint) string {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1<<16 - 1)
	return fmt.Sprintf("%04X%02X", random, id)
}
func ExecuteGoProgram(in string) (out string, err error) {
	response, err := http.PostForm("https://golang.google.cn/compile",
		url.Values{"body": []string{in}})
	if err != nil {
		return "", err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	out = string(body)
	compileErrors := gjson.Get(out, "compile_errors").String()
	output := gjson.Get(out, "output").String()
	if compileErrors != "" {
		out = "编译错误:" + compileErrors
	} else {
		out = output
	}
	return out, nil
}
func ExecuteJsProgram(in string) (out string, err error) {
	var halt = errors.New("block")
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Printf("Some code took to long! Stopping after: %v\n", duration)
				return
			}
			panic(caught) // Something else happened, repanic!
		}
	}()
	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)
	go func() {
		time.Sleep(2 * time.Second) // 只运行2秒，防止无限循环
		vm.Interrupt <- func() {
			panic(halt)
		}
	}()

	var logger string
	err = vm.Set("log", func(call otto.FunctionCall) otto.Value {
		outputs := make([]string, 0)
		for _, arg := range call.ArgumentList {
			outputs = append(outputs, arg.String())
		}
		logger = logger + strings.Join(outputs, " ") + "\n"
		return otto.Value{}
	})
	if err != nil {
		return "", err
	}
	in = "console.log=log;\n" + in //重载console.log为上一步的log函数,保存控制台输出
	_, err = vm.Run(in)
	if err != nil {
		return err.Error(), nil
	}
	return logger, nil
}
func ExecuteProgramSubject(language string, in string) (out string, err error) {
	language = strings.ToLower(language)
	if language == "javascript" || language == "js" {
		return ExecuteJsProgram(in)
	} else if language == "golang" || language == "go" {
		return ExecuteGoProgram(in)
	}
	return "", errors.New("unsupported language")
}
func SendEmail(to string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("Subject", "[辅助学习平台]")
	m.SetHeader("From", "2244363300@qq.com")
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "2244306600@qq.com", "Dan")抄送
	m.SetBody("text/html", fmt.Sprintf("你的注册验证码为<b>%s</b>，五分钟内有效", code))
	d := gomail.NewDialer("smtp.qq.com", 465, "2244363300@qq.com", "duxmplmmtfnedhha")
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
