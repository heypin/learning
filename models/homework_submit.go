package models

import "github.com/jinzhu/gorm"

type HomeworkSubmit struct {
	gorm.Model
	UserId            uint `json:"userId"`
	HomeworkPublishId uint `json:"homeworkPublishId"`
	TotalScore        uint `json:"totalScore"`
}
