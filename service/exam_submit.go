package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"learning/models"
	"time"
)

type ExamSubmitService struct {
	Id            uint
	UserId        uint
	ExamPublishId uint
	TotalScore    uint
	StartTime     time.Time
	FinishTime    time.Time
	Mark          *uint
	SubmitItems   []*models.ExamSubmitItem
}

func (s *ExamSubmitService) GetExamUserSubmitWithItems() (*models.ExamSubmit, error) {
	submit, err := models.GetExamUserSubmitWithItems(s.UserId, s.ExamPublishId)
	return submit, err
}
func (s *ExamSubmitService) CreateExamSubmit() (examSubmit *models.ExamSubmit, err error) {
	publish, err := models.GetExamPublishById(s.ExamPublishId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if now.Before(publish.BeginTime) || now.After(publish.EndTime) {
		return nil, errors.New("不在允许考试的时间范围内")
	}
	if submitRecord, err := models.GetUserExamSubmitRecord(s.UserId, s.ExamPublishId); err != nil {
		return nil, err
	} else if submitRecord != nil {
		return submitRecord, nil
	} else {
		submit := models.ExamSubmit{
			UserId:        s.UserId,
			ExamPublishId: s.ExamPublishId,
			StartTime:     time.Now(),
		}
		if id, err := models.AddExamSubmit(submit); err != nil {
			return nil, err
		} else {
			return models.GetExamSubmitById(id)
		}
	}
}
func (s *ExamSubmitService) GetExamSubmitById() (submit *models.ExamSubmit, err error) {
	return models.GetExamSubmitById(s.Id)
}
func (s *ExamSubmitService) GetExamSubmitsByPublishId() (submits []*models.ExamSubmit, err error) {
	submits, err = models.GetExamSubmitsByPublishId(s.ExamPublishId)
	return
}
func (s *ExamSubmitService) UpdateExamSubmitWithItems() error {
	submit := models.ExamSubmit{
		Model:       gorm.Model{ID: s.Id},
		Mark:        s.Mark,
		SubmitItems: s.SubmitItems,
	}
	err := models.UpdateExamSubmitWithItems(submit)
	return err
}
func (s *ExamSubmitService) UpdateExamSubmitById() error {
	submit := models.ExamSubmit{
		Model:      gorm.Model{ID: s.Id},
		StartTime:  s.StartTime,
		FinishTime: s.FinishTime,
		Mark:       s.Mark,
	}
	return models.UpdateExamSubmitById(submit)
}
