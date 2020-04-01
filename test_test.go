package main

import (
	"learning/conf"
	"learning/models"
	"log"
	"strings"
	"testing"
)

func TestMy(t *testing.T) {
	conf.SetUp()
	models.Setup()
	err := models.DeleteUserById(1)
	if err != nil {
		t.Fail()
	}
}

//models.GetAllComment()
//fmt.Println(len(comment))
//for _,v:=range comment{
//	fmt.Print(map[string]interface{}{
//		"id":v.ID,
//		"content":v.Content,
//	})
//	if len(v.Children)!=0{
//		fmt.Print("...child.....")
//		for _,c := range v.Children{
//			fmt.Print(map[string]interface{}{
//				"id":c.ID,
//				"content":c.Content,
//			})
//		}
//	}
//	fmt.Println()
//}
func TestCasbin(t *testing.T) {
	//authEnforcer, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	//authEnforcer.LoadPolicy()
	var a string = "a.b.jpg"
	s := strings.Split(a, ".")
	log.Println(s[len(s)-1])

}
