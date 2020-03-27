package models

import (
	"github.com/jinzhu/gorm"
)

const (
	ROLE_ADMIN = iota
	ROLE_STUDENT
	ROLE_TEACHER
)

type Student struct {
	gorm.Model
	Email    string
	Password string
	RealName string
	Number   uint
	Avatar   string
}

func AddStudent(s Student) (id uint, err error) {
	if err := db.Create(&s).Error; err != nil {
		return 0, err
	}
	return s.ID, nil
}
func GetStudentByEmail(email string) (*Student, error) {
	var s Student
	err := db.Where("email = ?", email).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
