package models

import "github.com/jinzhu/gorm"

type HomeworkLibItem struct {
	gorm.Model
	HomeworkLibId uint                     `json:"homeworkLibId"`
	Type          string                   `json:"type"`
	Question      string                   `json:"question"`
	Answer        string                   `json:"answer"`
	Score         uint                     `json:"score"`
	Options       []*HomeworkLibItemOption `json:"options" gorm:"foreignKey:homework_lib_item_id;association_foreignKey:id;"`
}

func AddLibItemAndOptions(item HomeworkLibItem) (id uint, err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&item).Error; err != nil { //创建时会关联创建选项
		tx.Rollback()
		return 0, err
	}
	if err := UpdateLibSubjectCountAndTotalScore(item, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return item.ID, nil
}

func UpdateLibSubjectCountAndTotalScore(item HomeworkLibItem, tx *gorm.DB) error {
	var subjectCount uint
	if err := tx.Model(&HomeworkLibItem{}).Where("homework_lib_id = ?", item.HomeworkLibId).
		Count(&subjectCount).Error; err != nil {
		return err
	}
	var totalScore uint
	if err := tx.Model(&HomeworkLibItem{}).Select("SUM(score)").
		Where("homework_lib_id = ?", item.HomeworkLibId).Row().Scan(&totalScore); err != nil {
		return err
	}
	lib := HomeworkLib{
		Model: gorm.Model{ID: item.HomeworkLibId},
	}
	if err := tx.Model(&lib).Updates(map[string]interface{}{
		"total_score":   totalScore,
		"subject_count": subjectCount,
	}).Error; err != nil {
		return err
	}
	return nil
}
func UpdateLibItemAndOptions(item HomeworkLibItem) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		//先删除选项再重新设置
		if item.Type == Subject_Single || item.Type == Subject_Multiple {
			if err := tx.Where("homework_lib_item_id = ?", item.ID).
				Delete(&HomeworkLibItemOption{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&item).Update(item).Error; err != nil {
			return err
		}
		if err := UpdateLibSubjectCountAndTotalScore(item, tx); err != nil {
			return err
		}
		return nil // 返回 nil 提交事务
	})
	return err
}
func GetLibItemsByLibId(libId uint) (items []*HomeworkLibItem, err error) {
	err = db.Where("homework_lib_id = ?", libId).
		Preload("Options").Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return items, nil
}

func DeleteLibItemAndOptions(id uint) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		//在删除前获取到,用于修改作业库的题目数和总分
		var item HomeworkLibItem
		if err := tx.Where("id = ?", id).First(&item).Error; err != nil {
			return err
		}
		//删除题目
		if err := tx.Where("id = ?", id).Delete(&HomeworkLibItem{}).Error; err != nil {
			return err
		}
		//删除选项
		if err := tx.Where("homework_lib_item_id = ?", id).
			Delete(&HomeworkLibItemOption{}).Error; err != nil {
			return err
		}

		if err := UpdateLibSubjectCountAndTotalScore(item, tx); err != nil {
			return err
		}
		return nil
	})
	return err
}
func GetHomeworkLibItemById(id uint) (*HomeworkLibItem, error) {
	var item HomeworkLibItem
	err := db.Where("id = ?", id).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, nil
}
