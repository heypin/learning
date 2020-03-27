package service

import (
	"learning/models"
	"learning/utils"
)

type StudentService struct {
	Email    string
	Password string
	RealName string
}

func (s *StudentService) Auth() bool {
	student, err := models.GetStudentByEmail(s.Email)
	if err == nil && student != nil && utils.CheckPassword(student.Password, s.Password) {
		return true
	}
	return false
}
func (s *StudentService) Register() (id uint, err error) {
	if _, err := models.GetStudentByEmail(s.Email); err != nil {
		return 0, err
	}
	student := models.Student{
		Email:    s.Email,
		Password: utils.Encrypt(s.Password),
		RealName: s.RealName,
	}
	if id, err := models.AddStudent(student); err != nil {
		return id, err
	}
	return id, nil
}

func (s *StudentService) GetStudentByEmail() (*models.Student, error) {
	student, err := models.GetStudentByEmail(s.Email)
	if err != nil {
		return nil, err
	}
	return student, nil
}
