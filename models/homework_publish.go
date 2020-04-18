package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type HomeworkPublish struct {
	gorm.Model
	ClassId       uint      `json:"classId"`
	HomeworkLibId uint      `json:"homeworkLibId"`
	BeginTime     time.Time `json:"beginTime"`
	EndTime       time.Time `json:"endTime"`
	Resubmit      uint      `json:"resubmit"`
}
