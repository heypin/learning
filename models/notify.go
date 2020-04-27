package models

import "github.com/jinzhu/gorm"

type Notify struct {
	gorm.Model
	CourseId uint   `json:"courseId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func AddNotify(n Notify) (id uint, err error) {
	if err := db.Create(&n).Error; err != nil {
		return 0, err
	}
	return n.ID, nil
}
func GetNotifyByCourseId(courseId uint) ([]*Notify, error) {
	var notifies []*Notify
	err := db.Where("course_id = ?", courseId).Find(&notifies).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return notifies, nil
}
func UpdateNotifyById(n Notify) (err error) {
	if err = db.Model(&n).Update(n).Error; err != nil {
		return err
	}
	return nil
}
func DeleteNotifyById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&Notify{}).Error; err != nil {
		return err
	}
	return nil
}
