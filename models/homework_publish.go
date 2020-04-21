package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type HomeworkPublish struct {
	gorm.Model
	ClassId       uint           `json:"classId"`
	HomeworkLibId uint           `json:"homeworkLibId"`
	BeginTime     time.Time      `json:"beginTime"`
	EndTime       time.Time      `json:"endTime"`
	Resubmit      *uint          `json:"resubmit"`
	HomeworkLib   HomeworkLib    `json:"homeworkLib" gorm:"foreignKey:id;association_foreignKey:homework_lib_id;"`
	SubmitRecord  HomeworkSubmit `json:"submitRecord" gorm:"foreignKey:homework_publish_id;association_foreignKey:id;"`
}

func AddHomeworkPublish(h HomeworkPublish) (id uint, err error) {
	if err := db.Create(&h).Error; err != nil {
		return 0, err
	}
	return h.ID, nil
}
func HasPublishHomework(classId uint, homeworkLibId uint) (bool, error) {
	var publish HomeworkPublish
	err := db.Where("class_id = ? AND homework_lib_id = ?", classId, homeworkLibId).
		First(&publish).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, nil
	}
}
func IsAllowResubmitHomework(publishId uint) (bool, error) {
	var publish HomeworkPublish
	err := db.Where("id = ?", publishId).First(&publish).Error
	if err != nil {
		return false, err
	}
	return *publish.Resubmit == Resubmit_Allow, nil
}
func UpdateHomeworkPublishById(h HomeworkPublish) error {
	if err := db.Model(&h).Update(h).Error; err != nil {
		return err
	}
	return nil
}
func GetHomeworkPublishesByClassId(classId uint) (publishes []*HomeworkPublish, err error) {
	err = db.Where("class_id = ?", classId).
		Preload("HomeworkLib").Find(&publishes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return publishes, nil
}
func GetHomeworkPublishesWithSubmitByClassId(classId uint, userId uint) (publishes []*HomeworkPublish, err error) {
	err = db.Where("class_id = ?", classId).Preload("HomeworkLib").
		Preload("SubmitRecord", "user_id = ?", userId).Find(&publishes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return publishes, nil
}
