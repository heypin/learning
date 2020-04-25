package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
	"time"
)

type HomeworkPublishService struct {
	Id            uint
	ClassId       uint
	HomeworkLibId uint
	BeginTime     time.Time
	EndTime       time.Time
	Resubmit      *uint
}

func (s *HomeworkPublishService) PublishHomework() (uint, error) {
	publish := models.HomeworkPublish{
		ClassId:       s.ClassId,
		HomeworkLibId: s.HomeworkLibId,
		BeginTime:     s.BeginTime,
		EndTime:       s.EndTime,
		Resubmit:      s.Resubmit,
	}
	if ok, err := models.HasPublishHomework(s.ClassId, s.HomeworkLibId); ok {
		return 0, nil
	} else if !ok && err == nil {
		if id, err := models.AddHomeworkPublish(publish); err == nil {
			return id, nil
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}
func (s *HomeworkPublishService) GetHomeworkPublishById() (publish *models.HomeworkPublish, err error) {
	return models.GetHomeworkPublishById(s.Id)
}
func (s *HomeworkPublishService) UpdateHomeworkPublishById() (err error) {
	publish := models.HomeworkPublish{
		Model:     gorm.Model{ID: s.Id},
		BeginTime: s.BeginTime,
		EndTime:   s.EndTime,
		Resubmit:  s.Resubmit,
	}
	err = models.UpdateHomeworkPublishById(publish)
	return
}
func (s *HomeworkPublishService) GetHomeworkPublishesByClassId() (publishes []*models.HomeworkPublish, err error) {
	publishes, err = models.GetHomeworkPublishesByClassId(s.ClassId)
	if err != nil {
		return publishes, err
	} else {
		for _, v := range publishes {
			v.SubmitCount, _ = models.GetHomeworkSubmitCountByPublishId(v.ID)
			v.UnMarkCount, _ = models.GetHomeworkUnmarkedCountByPublishId(v.ID)
		}
		return publishes, nil
	}
}
func (s *HomeworkPublishService) GetHomeworkPublishesWithSubmitByClassId(userId uint) (publishes []*models.HomeworkPublish, err error) {
	publishes, err = models.GetHomeworkPublishesWithSubmitByClassId(s.ClassId, userId)
	return
}
