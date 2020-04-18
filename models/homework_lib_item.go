package models

import "github.com/jinzhu/gorm"

type HomeworkLibItem struct {
	gorm.Model
	HomeworkLibId uint   `json:"homeworkLibId"`
	Type          string `json:"type"`
	Question      string `json:"question"`
	Answer        string `json:"answer"`
	Score         string `json:"score"`
}

func AddHomeworkLibItem(item HomeworkLibItem) (id uint, err error) {
	if err := db.Create(&item).Error; err != nil {
		return 0, err
	}
	return item.ID, nil
}
func UpdateHomeworkLibItem(item HomeworkLibItem) (err error) {
	if err := db.Model(&item).Update(item).Error; err != nil {
		return err
	}
	return nil
}
func GetLibItemsByLibId(libId uint) (items []*HomeworkLibItem, err error) {
	err = db.Where("homework_lib_id = ?", libId).Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return items, nil
}
func DeleteLibItemById(id uint) (err error) {
	if err = db.Where("id = ?", id).Delete(&HomeworkLibItem{}).Error; err != nil {
		return err
	}
	return nil
}
func DeleteLibItemsByLibId(libId uint) (err error) {
	if err = db.Where("homework_lib_id = ?", libId).Delete(&HomeworkLibItem{}).Error; err != nil {
		return err
	}
	return nil
}
