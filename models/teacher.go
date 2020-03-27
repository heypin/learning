package models

import "github.com/jinzhu/gorm"

type Teacher struct {
	gorm.Model
	Email    string
	Password string
	RealName string
	Avatar   string
}

func AddTeacher(t Teacher) (id uint, err error) {
	if err := db.Create(&t).Error; err != nil {
		return 0, err
	}
	return t.ID, nil
}
func GetTeacherByEmail(email string) (*Teacher, error) {
	var t Teacher
	err := db.Where("email = ?", email).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
