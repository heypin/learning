package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type HomeworkLibItemService struct {
	Id            uint
	HomeworkLibId uint
	Type          string
	Question      string
	Answer        *string
	Score         uint
	Options       []*models.HomeworkLibItemOption
}

func (s *HomeworkLibItemService) CreateLibItemAndOptions() (id uint, err error) {
	item := models.HomeworkLibItem{
		HomeworkLibId: s.HomeworkLibId,
		Type:          s.Type,
		Question:      s.Question,
		Answer:        s.Answer,
		Score:         s.Score,
		Options:       s.Options,
	}
	id, err = models.AddLibItemAndOptions(item)
	return
}
func (s *HomeworkLibItemService) UpdateLibItemAndOptions() (err error) {
	item := models.HomeworkLibItem{
		Model:         gorm.Model{ID: s.Id},
		HomeworkLibId: s.HomeworkLibId,
		Type:          s.Type,
		Question:      s.Question,
		Answer:        s.Answer,
		Score:         s.Score,
		Options:       s.Options,
	}
	err = models.UpdateLibItemAndOptions(item)
	return
}
func (s *HomeworkLibItemService) GetLibItemsByLibId() (items []*models.HomeworkLibItem, err error) {
	items, err = models.GetLibItemsByLibId(s.HomeworkLibId)
	return
}
func (s *HomeworkLibItemService) DeleteLibItemAndOptionsById() (err error) {
	err = models.DeleteLibItemAndOptions(s.Id)
	return
}
func (s *HomeworkLibItemService) GetHomeworkLibItemById() (item *models.HomeworkLibItem, err error) {
	item, err = models.GetHomeworkLibItemById(s.Id)
	return
}
