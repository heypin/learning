package service

import (
	"github.com/jinzhu/gorm"
	"learning/models"
)

type CommentService struct {
	Id          uint
	CourseId    uint
	ParentId    uint
	ReplyId     uint
	ReplyUserId uint
	UserId      uint
	Content     string
}

func (s *CommentService) GetCommentByCourseId() (c []*models.Comment, err error) {
	c, err = models.GetCommentByCourseId(s.CourseId)
	return
}
func (s *CommentService) GetCommentByUserId() (c []*models.Comment, err error) {
	c, err = models.GetCommentByUserId(s.UserId, s.CourseId)
	return
}
func (s *CommentService) GetCommentReplyToUser() (c []*models.Comment, err error) {
	c, err = models.GetCommentReplyToUser(s.UserId, s.CourseId)
	return
}
func (s *CommentService) AddComment() (id uint, err error) {
	comment := models.Comment{
		CourseId:    s.CourseId,
		ParentId:    s.ParentId,
		ReplyId:     s.ReplyId,
		ReplyUserId: s.ReplyUserId,
		UserId:      s.UserId,
		Content:     s.Content,
	}
	id, err = models.AddComment(comment)
	return
}

func (s *CommentService) UpdateCommentById() (err error) {
	comment := models.Comment{
		Model:   gorm.Model{ID: s.Id},
		Content: s.Content,
	}
	err = models.UpdateCommentById(comment)
	return
}
func (s *CommentService) DeleteCommentById() (err error) {
	err = models.DeleteCommentById(s.Id)
	return
}
