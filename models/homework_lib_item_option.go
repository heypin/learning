package models

type HomeworkLibItemOption struct {
	HomeworkLibItemId uint   `json:"homeworkLibItemId"`
	Sequence          string `json:"sequence"`
	Content           string `json:"content"`
}

func AddLibItemOption(option HomeworkLibItemOption) (err error) {

	if err = db.Create(&option).Error; err != nil {
		return err
	}
	return nil
}
func UpdateLibItemOption(option HomeworkLibItemOption) (err error) {
	if err := db.Model(&option).Update(option).Error; err != nil {
		return err
	}
	return nil
}
func DeleteOptionsByItemId(itemId uint) (err error) {
	if err = db.Where("homework_lib_item_id = ?", itemId).
		Delete(&HomeworkLibItemOption{}).Error; err != nil {
		return err
	}
	return nil
}
