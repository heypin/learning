package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/utils"
)

type ClassService struct {
	Id       uint
	CourseId uint
	//UserId    uint
	ClassName string
	ClassCode string
}

func (s *ClassService) CreateClass() (id uint, err error) {
	class := models.Class{
		CourseId:  s.CourseId,
		ClassName: s.ClassName,
	}
	if id, err = models.AddClass(class); err != nil {
		return 0, err
	}
	update := models.Class{
		Model:     gorm.Model{ID: id},
		ClassCode: utils.GenerateClassCode(id),
	}
	if err := models.UpdateClassById(update); err != nil {
		return 0, err
	}
	return id, nil
}
func (s *ClassService) GetClassByCourseId() (c []*models.Class, err error) {
	c, err = models.GetClassByCourseId(s.CourseId)
	return
}
