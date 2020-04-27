package main

import (
	"learning/cache"
	"learning/conf"
	"learning/models"
	"learning/routers"
	"log"
)

func main() {
	conf.SetUp()
	models.Setup()
	cache.SetUp()
	r := routers.InitRouters()
	if err := r.Run(":8080"); err != nil {
		log.Println("启动失败")
	}

}
