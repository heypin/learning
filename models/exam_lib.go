package models

import "github.com/jinzhu/gorm"

type ExamLib struct {
	gorm.Model
	CourseId     uint           `json:"courseId"`
	Name         string         `json:"name"`
	SubjectCount uint           `json:"subjectCount"`
	TotalScore   uint           `json:"totalScore"`
	Items        []*ExamLibItem `json:"items" gorm:"foreignKey:exam_lib_id;association_foreignKey:id;"`
}

func AddExamLib(lib ExamLib) (id uint, err error) {
	if err := db.Create(&lib).Error; err != nil {
		return 0, err
	}
	return lib.ID, nil
}
func UpdateExamLibById(lib ExamLib) error {
	if err := db.Model(&lib).Update(lib).Error; err != nil {
		return err
	}
	return nil
}
func GetExamLibWithItemsById(id uint) (*ExamLib, error) {
	var lib ExamLib
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
func GetExamLibsByCourseId(courseId uint) (libs []*ExamLib, err error) {
	err = db.Where("course_id = ?", courseId).Find(&libs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return libs, nil
}
