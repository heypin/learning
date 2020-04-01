package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	ClassId  int
	ParentId int
	ReplyId  int
	UserId   int
	Content  string
	Children []*Comment `gorm:"foreignkey:parent_id;association_foreignkey:id;"`
}

func GetAllComment() {
	var comment []Comment
	db.Where("parent_id IS NULL").
		Preload("Children").Find(&comment)

}
