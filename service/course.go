package service

import (
	"learning/models"
)

type CourseService struct {
	Id          uint
	UserId      uint
	Name        string
	Teacher     string
	Cover       string
	Description string
}

func (s *CourseService) AddCourse() (id uint, err error) {
	course := models.Course{
		UserId:      s.UserId,
		Name:        s.Name,
		Teacher:     s.Teacher,
		Cover:       s.Cover,
		Description: s.Description,
	}
	if id, err = models.AddCourse(course); err == nil {
		return id, nil
	}
	return 0, err
}
func (s *CourseService) GetCourseByUserId() ([]*models.Course, error) {
	courses, err := models.GetCourseByUserId(s.UserId)
	if err != nil {
		return nil, err
	}
	return courses, nil
}
func (s *CourseService) GetCourseById() (course *models.Course, err error) {
	return models.GetCourseById(s.Id)
}
