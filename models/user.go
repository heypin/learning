package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string   `json:"email"`
	Password string   `json:"-"`
	RealName string   `json:"realName"`
	Sex      uint     `json:"sex"`
	Number   string   `json:"number"`
	Avatar   string   `json:"avatar"`
	Classes  []*Class `json:"-" gorm:"many2many:class_member;"`
}

func AddUser(u User) (id uint, err error) {
	if err := db.Create(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}
func GetUserByEmail(email string) (*User, error) {
	var u User
	err := db.Where("email = ?", email).First(&u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, nil
}
func GetUserById(id uint) (*User, error) {
	var u User
	err := db.Where("id = ?", id).First(&u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, nil
}
func UpdateUserById(u User) (err error) {
	if err = db.Model(&u).Update(u).Error; err != nil {
		return err
	}
	return nil
}

//func DeleteUserById(id uint) (err error) {
//	if err = db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
