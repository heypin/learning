package main

import (
	"fmt"
	"learning/conf"
	"learning/models"
	"testing"
)

func TestMy(t *testing.T) {
	conf.SetUp()
	models.Setup()
	fmt.Println(models.GetFileById(19))
}

func TestCasbin(t *testing.T) {
	//authEnforcer, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	//authEnforcer.LoadPolicy()

}
