package models

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model
	UserId      uint   `json:"userId"`
	Name        string `json:"name"`
	Teacher     string `json:"teacher"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
}

func AddCourse(c Course) (id uint, err error) {
	if err := db.Create(&c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}
func GetCourseByUserId(id uint) ([]*Course, error) {
	var c []*Course
	err := db.Where("user_id = ?", id).Find(&c).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c, nil
}
func DeleteCourseById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&Course{}).Error; err != nil {
		return err
	}
	return nil
}
