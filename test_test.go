package main

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"learning/conf"
	"learning/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMy(t *testing.T) {
	conf.SetUp()
	models.Setup()
	//comments,_:=models.GetUserStudyClass(6)
	//fmt.Println(comments)
	//arr,_:=json.Marshal(comments)
	//log.Println(string(arr))
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
	in := `var a =1;log(a+5);` //"while(true){}"
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
		outputs := make([]string, 5)
		for _, arg := range call.ArgumentList {
			outputs = append(outputs, arg.String())
		}
		logger = logger + strings.Join(outputs, " ") + "\n"
		return otto.Value{}
	})
	_, err := vm.Run(in)
	fmt.Println(logger, err)
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
func TestM(t *testing.T) {
	v, err := strconv.Atoi("")
	log.Println(v, err)
}
