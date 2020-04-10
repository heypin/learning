package models

import "github.com/jinzhu/gorm"

type Chapter struct {
	gorm.Model
	UserId      uint    `json:"userId"`
	CourseId    uint    `json:"courseId"`
	ChapterName string  `json:"chapterName"`
	VideoName   *string `json:"videoName"`
}

func AddChapter(c Chapter) (id uint, err error) {
	if err := db.Create(&c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}
func GetChapterByCourseId(courseId uint) ([]*Chapter, error) {
	var chapters []*Chapter
	err := db.Where("course_id = ?", courseId).Find(&chapters).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return chapters, nil
}
func GetChapterById(id uint) (*Chapter, error) {
	var c Chapter
	err := db.Where("id = ?", id).First(&c).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &c, nil
}
func DeleteChapterById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&Chapter{}).Error; err != nil {
		return err
	}
	return nil
}
func UpdateChapterById(c Chapter) (err error) {

	if err = db.Model(&c).Update(c).Error; err != nil {
		return err
	}
	return nil
}
