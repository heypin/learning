package models

import "github.com/jinzhu/gorm"

type ExamLibItem struct {
	gorm.Model
	ExamLibId uint                 `json:"examLibId"`
	Type      string               `json:"type"`
	Question  string               `json:"question"`
	Answer    string               `json:"answer"`
	Score     uint                 `json:"score"`
	Options   []*ExamLibItemOption `json:"options" gorm:"foreignKey:exam_lib_item_id;association_foreignKey:id;"`
}

func UpdateExamLibSubjectCountAndTotalScore(item ExamLibItem, tx *gorm.DB) error {
	var subjectCount uint
	if err := tx.Model(&ExamLibItem{}).Where("exam_lib_id = ?", item.ExamLibId).
		Count(&subjectCount).Error; err != nil {
		return err
	}
	var totalScore uint
	if subjectCount != 0 {
		if err := tx.Model(&ExamLibItem{}).Select("SUM(score)").
			Where("exam_lib_id = ?", item.ExamLibId).Row().Scan(&totalScore); err != nil {
			return err
		}
	}
	lib := ExamLib{
		Model: gorm.Model{ID: item.ExamLibId},
	}
	if err := tx.Model(&lib).Updates(map[string]interface{}{
		"total_score":   totalScore,
		"subject_count": subjectCount,
	}).Error; err != nil {
		return err
	}
	return nil
}
func AddExamLibItemAndOptions(item ExamLibItem) (id uint, err error) {
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
	if err := UpdateExamLibSubjectCountAndTotalScore(item, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return item.ID, nil
}
func GetExamLibItemsByLibId(libId uint) (items []*ExamLibItem, err error) {
	err = db.Where("exam_lib_id = ?", libId).
		Preload("Options").Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return items, nil
}
func UpdateExamLibItemAndOptions(item ExamLibItem) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		//先删除选项再重新设置
		if item.Type == SubjectSingle || item.Type == SubjectMultiple {
			if err := tx.Where("exam_lib_item_id = ?", item.ID).
				Delete(&ExamLibItemOption{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&item).Update(item).Error; err != nil {
			return err
		}
		if err := UpdateExamLibSubjectCountAndTotalScore(item, tx); err != nil {
			return err
		}
		return nil // 返回 nil 提交事务
	})
	return err
}
func DeleteExamLibItemAndOptions(id uint) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		//在删除前获取到,用于修改题库的题目数和总分
		var item ExamLibItem
		if err := tx.Where("id = ?", id).First(&item).Error; err != nil {
			return err
		}
		//删除题目
		if err := tx.Where("id = ?", id).Delete(&ExamLibItem{}).Error; err != nil {
			return err
		}
		//删除选项
		if err := tx.Where("exam_lib_item_id = ?", id).
			Delete(&ExamLibItemOption{}).Error; err != nil {
			return err
		}
		if err := UpdateExamLibSubjectCountAndTotalScore(item, tx); err != nil {
			return err
		}
		return nil
	})
	return err
}
func GetExamLibItemById(id uint) (*ExamLibItem, error) {
	var item ExamLibItem
	err := db.Where("id = ?", id).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, nil
}
