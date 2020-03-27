package main

import (
	"learning/conf"
	"learning/models"
	"learning/routers"
	"log"
)

func main() {
	conf.SetUp()
	models.Setup()
	r := routers.InitRouters()
	if err := r.Run(":8080"); err != nil {
		log.Println("启动失败")
	}

}
