package models

import "github.com/jinzhu/gorm"

type Class struct {
	gorm.Model
	CourseId  uint   `json:"courseId"`
	UserId    uint   `json:"userId"`
	ClassName string `json:"className"`
	ClassCode string `json:"classCode"`
}

func AddClass(c Class) (id uint, err error) {
	if err := db.Create(&c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}
func UpdateClassById(c Class) (err error) {
	if err := db.Model(&c).Update(c).Error; err != nil {
		return err
	}
	return nil
}
func GetClassByCourseId(id uint) ([]*Class, error) {
	var c []*Class
	err := db.Where("course_id = ?", id).Find(&c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}
func DeleteClassById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&Class{}).Error; err != nil {
		return err
	}
	return nil
}
