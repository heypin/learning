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
	vm := otto.New()
	value, err := vm.Run(in)
	if err != nil {
		return "", err
	}
	log.Println(value)
	return value.String(), nil
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
