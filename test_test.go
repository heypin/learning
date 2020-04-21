package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"learning/conf"
	"learning/models"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestMy(t *testing.T) {
	conf.SetUp()
	models.Setup()
	//comments,_:=models.GetUserStudyClass(6)
	//fmt.Println(comments)
	//arr,_:=json.Marshal(comments)
	//log.Println(string(arr))
}
func TestJsProgram(t *testing.T) {
	in := "console.log(12);"
	vm := otto.New()
	value, _ := vm.Run(in)
	out := value.
		log.Println(out)
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
	str := "language:java;"
	fmt.Println(len(strings.SplitN(str, ";", 2)))
	fmt.Println(len(strings.SplitAfter(str, ";")))
	fmt.Println(len(strings.SplitAfterN(str, ";", 2)))
}
