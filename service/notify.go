package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type NotifyService struct {
	Id       uint
	UserId   uint
	CourseId uint
	Title    string
	Content  string
}

func (s *NotifyService) AddNotify() (id uint, err error) {
	n := models.Notify{
		UserId:   s.UserId,
		CourseId: s.CourseId,
		Title:    s.Title,
		Content:  s.Content,
	}
	id, err = models.AddNotify(n)
	return
}
func (s *NotifyService) GetNotifyByCourseId() (n []*models.Notify, err error) {
	n, err = models.GetNotifyByCourseId(s.CourseId)
	return
}
func (s *NotifyService) UpdateNotifyById() (err error) {
	n := models.Notify{
		Model:   gorm.Model{ID: s.Id},
		Title:   s.Title,
		Content: s.Content,
	}
	err = models.UpdateNotifyById(n)
	return
}
func (s *NotifyService) DeleteNotifyById() (err error) {
	err = models.DeleteNotifyById(s.Id)
	return
}
