package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
	"time"
)

type ExamPublishService struct {
	Id        uint
	ClassId   uint
	ExamLibId uint
	BeginTime time.Time
	EndTime   time.Time
	Duration  uint
}

func (s *ExamPublishService) PublishExam() (uint, error) {
	publish := models.ExamPublish{
		ClassId:   s.ClassId,
		ExamLibId: s.ExamLibId,
		BeginTime: s.BeginTime,
		EndTime:   s.EndTime,
		Duration:  s.Duration,
	}
	if ok, err := models.HasPublishExam(s.ClassId, s.ExamLibId); ok {
		return 0, nil
	} else if !ok && err == nil {
		if id, err := models.AddExamPublish(publish); err == nil {
			return id, nil
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}
func (s *ExamPublishService) GetExamPublishById() (publish *models.ExamPublish, err error) {
	return models.GetExamPublishById(s.Id)
}
func (s *ExamPublishService) GetExamPublishesByClassId() (publishes []*models.ExamPublish, err error) {
	publishes, err = models.GetExamPublishesByClassId(s.ClassId)
	if err != nil {
		return publishes, err
	} else {
		for _, v := range publishes {
			v.SubmitCount, _ = models.GetExamSubmitCountByPublishId(v.ID)
			v.UnMarkCount, _ = models.GetExamUnmarkedCountByPublishId(v.ID)
		}
		return publishes, nil
	}
}
func (s *ExamPublishService) UpdateExamPublishById() (err error) {
	publish := models.ExamPublish{
		Model:     gorm.Model{ID: s.Id},
		BeginTime: s.BeginTime,
		EndTime:   s.EndTime,
		Duration:  s.Duration,
	}
	err = models.UpdateExamPublishById(publish)
	return
}
func (s *ExamPublishService) GetExamPublishesWithUserSubmitByClassId(userId uint) (publishes []*models.ExamPublish, err error) {
	publishes, err = models.GetExamPublishesWithUserSubmitByClassId(s.ClassId, userId)
	return
}
