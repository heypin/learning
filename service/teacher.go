package service

import (
	"learning/models"
	"learning/utils"
)

type TeacherService struct {
	Email    string
	Password string
	RealName string
}

func (t *TeacherService) Auth() bool {
	teacher, err := models.GetTeacherByEmail(t.Email)
	if err == nil && teacher != nil && teacher.Password == utils.Encrypt(t.Password) {
		return true
	}
	return false
}
func (t *TeacherService) Register() (id uint, err error) {
	if _, err := models.GetTeacherByEmail(t.Email); err != nil {
		return 0, err
	}
	teacher := models.Teacher{
		Email:    t.Email,
		Password: utils.Encrypt(t.Password),
		RealName: t.RealName,
	}
	if id, err := models.AddTeacher(teacher); err != nil {
		return id, err
	}
	return id, nil
}
