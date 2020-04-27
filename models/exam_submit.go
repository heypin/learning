package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ExamSubmit struct {
	gorm.Model
	UserId        uint              `json:"userId"`
	ExamPublishId uint              `json:"examPublishId"`
	TotalScore    uint              `json:"totalScore"`
	StartTime     time.Time         `json:"startTime"`
	FinishTime    *time.Time        `json:"finishTime"`
	Mark          *uint             `json:"mark"`
	SubmitItems   []*ExamSubmitItem `json:"submitItems" gorm:"foreignKey:exam_submit_id;association_foreignKey:id;"`
	User          User              `json:"user" gorm:"foreignKey:id;association_foreignKey:user_id;"`
}

func GetExamUserSubmitWithItems(userId uint, publishId uint) (*ExamSubmit, error) {
	var submit ExamSubmit
	err := db.Where("user_id = ? AND exam_publish_id = ?", userId, publishId).
		Preload("SubmitItems").First(&submit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &submit, nil
}
func AddExamSubmit(submit ExamSubmit) (id uint, err error) {
	if err := db.Create(&submit).Error; err != nil {
		return 0, err
	}
	return submit.ID, nil
}

//有记录代表用户开始了考试
func GetUserExamSubmitRecord(userId uint, publishId uint) (*ExamSubmit, error) {
	var submit ExamSubmit
	err := db.Where("user_id = ? AND exam_publish_id = ?", userId, publishId).
		First(&submit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &submit, nil
}
func GetExamSubmitById(id uint) (*ExamSubmit, error) {
	var submit ExamSubmit
	err := db.Where("id = ?", id).First(&submit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &submit, nil
}

//获取考试所有提交记录
func GetExamSubmitsByPublishId(publishId uint) (submits []*ExamSubmit, err error) {
	err = db.Where("exam_publish_id = ?", publishId).
		Preload("User").Find(&submits).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return submits, nil
}

//获取考试提交数量
func GetExamSubmitCountByPublishId(publishId uint) (count uint, err error) {
	err = db.Model(&ExamSubmit{}).Where("exam_publish_id = ?", publishId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

//获取考试提交但未批改的数量
func GetExamUnmarkedCountByPublishId(publishId uint) (count uint, err error) {
	err = db.Model(&ExamSubmit{}).Where("mark = 0 AND exam_publish_id = ?", publishId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func CalculateUserExamScore(submit ExamSubmit, tx *gorm.DB) (err error) {
	var totalScore uint
	if err = tx.Model(&ExamSubmitItem{}).Select("SUM(score)").
		Where("exam_submit_id = ?", submit.ID).Row().Scan(&totalScore); err != nil {
		return err
	}
	if err = tx.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Model(&submit).Update("total_score", totalScore).Error; err != nil {
		return err
	}
	return nil
}
func UpdateExamSubmitById(submit ExamSubmit) error {
	if err := db.Model(&submit).Update(submit).Error; err != nil {
		return err
	}
	return nil
}
func UpdateExamSubmitWithItems(submit ExamSubmit) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range submit.SubmitItems {
			item.ExamSubmitId = submit.ID
			if item.ID == 0 { //有ID更新它，无主键创建它
				if err := tx.Create(item).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(item).Update(item).Error; err != nil {
					return err
				}
			}
		}
		submit.SubmitItems = nil //子项已在上一步更新，gorm关联更新会更新关联结构体零值,这里置为空,或者使用
		//Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).
		if err := tx.Model(&submit).Updates(submit).Error; err != nil {
			return err
		}
		if err := CalculateUserExamScore(submit, tx); err != nil {
			return err
		}
		return nil
	})
	return err
}
