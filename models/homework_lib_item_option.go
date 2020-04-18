package models

import "github.com/jinzhu/gorm"

type HomeworkLibItemOption struct {
	gorm.Model
	HomeworkLibItemId uint   `json:"homeworkLibItemId"`
	Sequence          string `json:"sequence"`
	Content           string `json:"content"`
}
