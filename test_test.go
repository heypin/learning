package main

import (
	"learning/conf"
	"learning/models"
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
