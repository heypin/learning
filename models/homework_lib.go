package models

import "github.com/jinzhu/gorm"

type HomeworkLib struct {
	gorm.Model
	CourseId     uint               `json:"courseId"`
	Name         string             `json:"name"`
	SubjectCount uint               `json:"subjectCount"`
	TotalScore   uint               `json:"totalScore"`
	Items        []*HomeworkLibItem `json:"items" gorm:"foreignKey:homework_lib_id;association_foreignKey:id;"`
}

func AddHomeworkLib(h HomeworkLib) (id uint, err error) {
	if err := db.Create(&h).Error; err != nil {
		return 0, err
	}
	return h.ID, nil
}
func UpdateHomeworkLibById(h HomeworkLib) error {
	if err := db.Model(&h).Update(h).Error; err != nil {
		return err
	}
	return nil
}

//获取作业库及其库中的题和每到题
func GetHomeworkLibWithItemsById(id uint) (*HomeworkLib, error) {
	var lib HomeworkLib
	err := db.Where("id = ?", id).Preload("Items").
		Preload("Items.Options").First(&lib).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &lib, nil
}
func GetHomeworkLibsByCourseId(courseId uint) (libs []*HomeworkLib, err error) {
	err = db.Where("course_id = ?", courseId).Find(&libs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return libs, nil
}
