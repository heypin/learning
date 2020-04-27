package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ExamPublish struct {
	gorm.Model
	ClassId      uint       `json:"classId"`
	ExamLibId    uint       `json:"examLibId"`
	BeginTime    time.Time  `json:"beginTime"`
	EndTime      time.Time  `json:"endTime"`
	Duration     uint       `json:"duration"`
	ExamLib      ExamLib    `json:"examLib" gorm:"foreignKey:id;association_foreignKey:Exam_lib_id;"`
	SubmitRecord ExamSubmit `json:"submitRecord" gorm:"foreignKey:exam_publish_id;association_foreignKey:id;"` //保存某个用户的提交记录
	SubmitCount  uint       `json:"submitCount" gorm:"-"`                                                      //提交数
	UnMarkCount  uint       `json:"unMarkCount" gorm:"-"`                                                      //未批阅数
}

func AddExamPublish(publish ExamPublish) (id uint, err error) {
	if err := db.Create(&publish).Error; err != nil {
		return 0, err
	}
	return publish.ID, nil
}
func HasPublishExam(classId uint, examLibId uint) (bool, error) {
	var publish ExamPublish
	err := db.Where("class_id = ? AND exam_lib_id = ?", classId, examLibId).
		First(&publish).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, nil
	}
}
func UpdateExamPublishById(publish ExamPublish) error {
	if err := db.Model(&publish).Update(publish).Error; err != nil {
		return err
	}
	return nil
}
func GetExamPublishesByClassId(classId uint) (publishes []*ExamPublish, err error) {
	err = db.Where("class_id = ?", classId).
		Preload("ExamLib").Find(&publishes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return publishes, nil
}
func GetExamPublishById(id uint) (*ExamPublish, error) {
	var publish ExamPublish
	err := db.Where("id = ?", id).Preload("ExamLib").First(&publish).Error
	if err != nil {
		return nil, err
	}
	return &publish, nil
}
func GetExamPublishesWithUserSubmitByClassId(classId uint, userId uint) (publishes []*ExamPublish, err error) {
	err = db.Where("class_id = ?", classId).Preload("ExamLib").
		Preload("SubmitRecord", "user_id = ?", userId).Find(&publishes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return publishes, nil
}
