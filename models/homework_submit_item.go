package models

import "github.com/jinzhu/gorm"

type HomeworkSubmitItem struct {
	gorm.Model
	HomeworkSubmitId  uint   `json:"homeworkSubmitId"`
	HomeworkLibItemId uint   `json:"homeworkLibItemId"`
	Answer            string `json:"answer"`
	Score             *uint  `json:"score"`
}
