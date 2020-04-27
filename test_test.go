package main

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis/v7"
	"github.com/robertkrimen/otto"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"learning/conf"
	"learning/models"
	"learning/utils"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMy(t *testing.T) {
	conf.SetUp()
	models.Setup()
	//submits, _ := models.GetExamSubmitsByPublishId(1)

}
func TestEncrypt(t *testing.T) {
	hashed := utils.Encrypt("12345678")
	fmt.Println(hashed)
}

var halt = errors.New("block")

func TestJsProgram(t *testing.T) {
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
		fmt.Printf("Ran code successfully: %v\n", duration)
	}()
	in := `console.log=log;console.log(1+2);` //while(true){}
	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)
	go func() {
		time.Sleep(2 * time.Second) // Stop after two seconds，防止无限循环
		vm.Interrupt <- func() {
			panic(halt)
		}
	}()

	var logger string
	vm.Set("log", func(call otto.FunctionCall) otto.Value {
		outputs := make([]string, 0)
		for _, arg := range call.ArgumentList {
			outputs = append(outputs, arg.String())
		}
		logger = logger + strings.Join(outputs, " ") + "\n"
		return otto.Value{}
	})
	vm.Run(in)
	fmt.Println("logger", logger, "12")
}

func TestGoProgram(t *testing.T) {
	str := `package main
	import "fmt"
	func main() {
		fmt.Println("Hello, 世界")
	}`
	resp, err := http.Post("https://golang.google.cn/compile",
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(fmt.Sprintf("body=%s", str)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(body))
}
func TestGo(t *testing.T) {
	var answer = "B,A"
	var right = "A,B"
	answerSet := mapset.NewSet()
	for _, v := range strings.Split(answer, ",") {
		answerSet.Add(v)
	}
	rightSet := mapset.NewSet()
	for _, v := range strings.Split(right, ",") {
		rightSet.Add(v)
	}
	fmt.Println(rightSet.Equal(answerSet))
}
func TestWriteExcel(t *testing.T) {

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	title := map[string]string{
		"A1": "帐号",
		"B1": "姓名",
		"C1": "学号",
		"D1": "总分",
		"E1": "开考时间",
		"F1": "完成时间",
	}
	for k, v := range title {
		_ = f.SetCellValue("Sheet1", k, v)
	}
	err := f.SetColWidth("Sheet1", "A", "F", 20)
	log.Println(err)

	if err := f.SaveAs("../a.xlsx"); err != nil {
		fmt.Println(err)
	}
}

//duxmplmmtfnedhha  QQ邮箱授权码
func TestSendMail(t *testing.T) {

	m := gomail.NewMessage()
	m.SetHeader("Subject", "[辅助学习平台]")
	m.SetHeader("From", "2244363300@qq.com")
	m.SetHeader("To", "2244306600@qq.com")
	//m.SetAddressHeader("Cc", "2244306600@qq.com", "Dan")抄送
	m.SetBody("text/html", fmt.Sprintf("你的注册验证码为<b>%s</b>，五分钟内有效", "1234"))
	d := gomail.NewDialer("smtp.qq.com", 465, "2244363300@qq.com", "duxmplmmtfnedhha")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}
func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.AppConfig.Redis.Addr,
		Password: conf.AppConfig.Redis.Password,
		DB:       0, // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}
