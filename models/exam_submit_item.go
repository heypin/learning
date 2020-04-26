package models

import "github.com/jinzhu/gorm"

type ExamSubmitItem struct {
	gorm.Model
	ExamSubmitId  uint   `json:"examSubmitId"`
	ExamLibItemId uint   `json:"examLibItemId"`
	Answer        string `json:"answer"`
	Score         *uint  `json:"score"`
}
