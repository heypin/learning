package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"learning/conf"
	"log"
)

var db *gorm.DB

func Setup() {
	var err error
	db, err = gorm.Open("mysql", conf.AppConfig.DB.Url)

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	db.SingularTable(true)
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return setting.DatabaseSetting.TablePrefix + defaultTableName
	//}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
