package models

import "github.com/jinzhu/gorm"

type ClassMember struct {
	gorm.Model
	ClassId uint
	UserId  uint
	Status  uint
}
