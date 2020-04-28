package models

import "github.com/jinzhu/gorm"

type Class struct {
	gorm.Model
	CourseId  uint    `json:"courseId"`
	ClassName string  `json:"className"`
	ClassCode string  `json:"classCode"`
	Course    Course  `json:"course" gorm:"foreignKey:id;association_foreignKey:course_id;"`
	Users     []*User `json:"-" gorm:"many2many:class_member;"`
}

func AddClass(c Class) (id uint, err error) {
	if err := db.Create(&c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}
func UpdateClassById(c Class) (err error) {
	if err = db.Model(&c).Update(c).Error; err != nil {
		return err
	}
	return nil
}
func GetClassByCourseId(id uint) ([]*Class, error) {
	var c []*Class
	err := db.Where("course_id = ?", id).Find(&c).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c, nil
}
func GetClassByClassCode(code string) (*Class, error) {
	var class Class
	err := db.Where("class_code = ?", code).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

//func DeleteClassById(id uint) (err error) {
//	if err = db.Where("id = ?", id).Delete(&Class{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
