package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type ChapterService struct {
	Id          uint
	UserId      uint
	CourseId    uint
	ChapterName string
	VideoName   *string
}

func (s *ChapterService) AddChapter() (id uint, err error) {
	c := models.Chapter{
		UserId:      s.UserId,
		CourseId:    s.CourseId,
		ChapterName: s.ChapterName,
		VideoName:   s.VideoName,
	}
	id, err = models.AddChapter(c)
	return
}
func (s *ChapterService) GetChapterByCourseId() (c []*models.Chapter, err error) {
	c, err = models.GetChapterByCourseId(s.CourseId)
	return
}
func (s *ChapterService) DeleteChapterById() (err error) {
	err = models.DeleteChapterById(s.Id)
	return
}
func (s *ChapterService) GetChapterById() (c *models.Chapter, err error) {
	c, err = models.GetChapterById(s.Id)
	return
}
func (s *ChapterService) UpdateChapterById() (err error) {
	c := models.Chapter{
		Model:       gorm.Model{ID: s.Id},
		UserId:      s.UserId,
		CourseId:    s.CourseId,
		ChapterName: s.ChapterName,
		VideoName:   s.VideoName,
	}
	err = models.UpdateChapterById(c)
	return
}
