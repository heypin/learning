package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	CourseId    uint       `json:"courseId"`
	ParentId    uint       `json:"parentId" gorm:"default:0"`
	ReplyId     uint       `json:"replyId" gorm:"default:0"`
	ReplyUserId uint       `json:"replyUserId" gorm:"default:0"`
	UserId      uint       `json:"userId"`
	User        User       `json:"user" gorm:"foreignkey:id;association_foreignkey:user_id;"`
	Content     string     `json:"content"`
	Children    []*Comment `json:"children" gorm:"foreignkey:parent_id;association_foreignkey:id;"`
}

func GetCommentByCourseId(courseId uint) ([]*Comment, error) {
	var comments []*Comment
	err := db.Where("parent_id = 0 AND course_id = ?", courseId).
		Preload("User").Preload("Children").
		Preload("Children.User").Find(&comments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return comments, nil
}
func GetCommentByUserId(userId uint, courseId uint) ([]*Comment, error) {
	var comments []*Comment
	err := db.Where("user_id = ?", userId).
		Where("course_id = ?", courseId).Preload("User").
		Find(&comments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return comments, nil
}
func GetCommentReplyToUser(userId uint, courseId uint) ([]*Comment, error) {
	var comments []*Comment
	err := db.Where("reply_user_id = ?", userId).
		Where("course_id = ?", courseId).Preload("User").
		Find(&comments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return comments, nil
}
func AddComment(c Comment) (id uint, err error) {
	if err := db.Create(&c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}
func UpdateCommentById(c Comment) (err error) {
	if err = db.Model(&c).Update(c).Error; err != nil {
		return err
	}
	return nil
}
func DeleteCommentById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}
