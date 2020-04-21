package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type HomeworkLibService struct {
	Id           uint
	CourseId     uint
	Name         string
	SubjectCount uint
	TotalScore   uint
}

func (s *HomeworkLibService) AddHomeworkLib() (id uint, err error) {
	homeworkLib := models.HomeworkLib{
		CourseId: s.CourseId,
		Name:     s.Name,
	}
	id, err = models.AddHomeworkLib(homeworkLib)
	return
}
func (s *HomeworkLibService) UpdateHomeworkLibById() (err error) {
	homeworkLib := models.HomeworkLib{
		Model:        gorm.Model{ID: s.Id},
		Name:         s.Name,
		SubjectCount: s.SubjectCount,
		TotalScore:   s.TotalScore,
	}
	err = models.UpdateHomeworkLibById(homeworkLib)
	return
}
func (s *HomeworkLibService) GetHomeworkLibWithItemsById() (lib *models.HomeworkLib, err error) {
	lib, err = models.GetHomeworkLibWithItemsById(s.Id)
	return
}
func (s *HomeworkLibService) GetHomeworkLibsByCourseId() (libs []*models.HomeworkLib, err error) {
	libs, err = models.GetHomeworkLibsByCourseId(s.CourseId)
	return
}
