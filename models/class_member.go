package models

import (
	"github.com/jinzhu/gorm"
)

type ClassMember struct {
	ClassId uint
	UserId  uint
}

func GetUsersByClassId(classId uint) (users []*User, err error) {
	class := Class{
		Model: gorm.Model{ID: classId},
	}
	err = db.Model(&class).Related(&users, "Users").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return users, nil
}
func GetClassesByUserId(userId uint) (classes []*Class, err error) {
	var user User
	err = db.Where("id = ?", userId).Preload("Classes").
		Preload("Classes.Course").First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return user.Classes, nil
}

func ExistClassMemberRecord(userId uint, classId uint) (bool, error) {
	var classMember ClassMember
	err := db.Where("user_id = ? AND class_id = ?", userId, classId).
		First(&classMember).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, nil
	}
}
func AddClassMember(userId uint, classId uint) error {
	if err := db.Create(&ClassMember{
		ClassId: classId,
		UserId:  userId,
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteClassMember(userId uint, classId uint) (err error) {
	if err = db.Where("user_id = ? AND class_id = ?", userId, classId).
		Delete(&ClassMember{}).Error; err != nil {
		return err
	}
	return nil
}

//func GetClassesByUserId(userId uint)(classes []*Class,err error){
//	user := User{
//		Model:    gorm.Model{ID:userId},
//	}
//	err = db.Model(&user).Related(&classes,"Classes").Error
//	if err!=nil && err != gorm.ErrRecordNotFound{
//		return nil,err
//	}
//	return classes,nil
//}

//select * from course where id in (
//	select course_id from class where id in (
//		select class_id from class_member WHERE user_id=6))

//func GetUserStudyCourse(userId uint)(courses []*Course,err error){
//	err = db.Where("id IN (?)",
//		db.Table("class").Select("course_id").Where("id in (?)",
//			db.Table("class_member").Select("class_id").Where("user_id = ?",
//				userId).SubQuery()).SubQuery()).Find(&courses).Error
//	if err!=nil && err!=gorm.ErrRecordNotFound{
//		return nil,err
//	}
//	return courses,nil
//}
