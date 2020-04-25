package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type ExamLibService struct {
	Id           uint
	CourseId     uint
	Name         string
	SubjectCount uint
	TotalScore   uint
}

func (s *ExamLibService) AddExamLib() (id uint, err error) {
	examLib := models.ExamLib{
		CourseId: s.CourseId,
		Name:     s.Name,
	}
	id, err = models.AddExamLib(examLib)
	return
}
func (s *ExamLibService) UpdateExamLibById() error {
	examLib := models.ExamLib{
		Model:        gorm.Model{ID: s.Id},
		Name:         s.Name,
		SubjectCount: s.SubjectCount,
		TotalScore:   s.TotalScore,
	}
	err := models.UpdateExamLibById(examLib)
	return err
}
func (s *ExamLibService) GetExamLibWithItemsById() (lib *models.ExamLib, err error) {
	lib, err = models.GetExamLibWithItemsById(s.Id)
	return
}
func (s *ExamLibService) GetExamLibsByCourseId() (libs []*models.ExamLib, err error) {
	libs, err = models.GetExamLibsByCourseId(s.CourseId)
	return
}
