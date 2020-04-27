package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	UserId        uint   `json:"userId"`
	ParentId      uint   `json:"parentId" gorm:"default:0"`
	CourseId      uint   `json:"courseId"`
	Filename      string `json:"filename"`
	LocalFilepath string `json:"localFilepath"`
	LocalFilename string `json:"localFilename"`
	Size          uint   `json:"size"`
}

func AddFile(f File) (id uint, err error) {
	if err := db.Create(&f).Error; err != nil {
		return 0, err
	}
	return f.ID, nil
}
func GetChildFileByCourseId(courseId uint, parentId uint) (f []*File, err error) {
	err = db.Where("course_id = ?", courseId).
		Where("parent_id = ?", parentId).Find(&f).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return f, nil
}
func GetFileById(id uint) (*File, error) {
	var f File
	err := db.Where("id = ?", id).First(&f).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &f, nil
}

func DeleteFileById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}

//func UpdateFileById(f File) (err error) {
//	if err := db.Model(&f).Update(f).Error; err != nil {
//		return err
//	}
//	return nil
//}
