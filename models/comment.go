package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	CourseId int        `json:"courseId"`
	ParentId int        `json:"parentId"`
	ReplyId  int        `json:"replyId"`
	UserId   int        `json:"userId"`
	Content  string     `json:"content"`
	Children []*Comment `json:"children" gorm:"foreignkey:parent_id;association_foreignkey:id;"`
}

func GetAllComment() []Comment {
	var comment []Comment
	db.Where("parent_id IS NULL").
		Preload("Children").Find(&comment)
	return comment
}
