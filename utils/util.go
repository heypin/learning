package utils

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"io/ioutil"
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
