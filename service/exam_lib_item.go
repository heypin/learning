package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type ExamLibItemService struct {
	Id        uint
	ExamLibId uint
	Type      string
	Question  string
	Answer    string
	Score     uint
	Options   []*models.ExamLibItemOption
}

func (s *ExamLibItemService) CreateExamLibItemAndOptions() (id uint, err error) {
	item := models.ExamLibItem{
		ExamLibId: s.ExamLibId,
		Type:      s.Type,
		Question:  s.Question,
		Answer:    s.Answer,
		Score:     s.Score,
		Options:   s.Options,
	}
	id, err = models.AddExamLibItemAndOptions(item)
	return
}
func (s *ExamLibItemService) UpdateExamLibItemAndOptions() (err error) {
	item := models.ExamLibItem{
		Model:    gorm.Model{ID: s.Id},
		Type:     s.Type,
		Question: s.Question,
		Answer:   s.Answer,
		Score:    s.Score,
		Options:  s.Options,
	}
	err = models.UpdateExamLibItemAndOptions(item)
	return
}
func (s *ExamLibItemService) GetExamLibItemsByLibId() (items []*models.ExamLibItem, err error) {
	items, err = models.GetExamLibItemsByLibId(s.ExamLibId)
	return
}
func (s *ExamLibItemService) DeleteExamLibItemAndOptions() (err error) {
	err = models.DeleteExamLibItemAndOptions(s.Id)
	return
}
func (s *ExamLibItemService) GetExamLibItemById() (item *models.ExamLibItem, err error) {
	item, err = models.GetExamLibItemById(s.Id)
	return
}
