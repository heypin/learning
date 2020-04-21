package models

import "github.com/jinzhu/gorm"

type HomeworkSubmit struct {
	gorm.Model
	UserId            uint                  `json:"userId"`
	HomeworkPublishId uint                  `json:"homeworkPublishId"`
	TotalScore        uint                  `json:"totalScore"`
	Mark              *uint                 `json:"mark"`
	SubmitItems       []*HomeworkSubmitItem `json:"submitItems" gorm:"foreignKey:homework_submit_id;association_foreignKey:id;"`
}

func GetHomeworkSubmitWithItems(userId uint, publishId uint) (*HomeworkSubmit, error) {
	var submit HomeworkSubmit
	err := db.Where("user_id = ? AND homework_publish_id = ?", userId, publishId).
		Preload("SubmitItems").First(&submit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &submit, nil
}

//获取作业所有提交记录
func GetHomeworkSubmitsByPublishId(publishId uint) (submits []*HomeworkSubmit, err error) {
	err = db.Where("homework_publish_id = ?", publishId).Find(&submits).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return submits, nil
}

//获取作业提交数量
func GetHomeworkSubmitCountByPublishId(publishId uint) (count uint, err error) {
	err = db.Model(&HomeworkSubmit{}).Where("homework_publish_id = ?", publishId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

//获取作业提交但未批改的数量
func GetHomeworkUnmarkedCountByPublishId(publishId uint) (count uint, err error) {
	err = db.Model(&HomeworkSubmit{}).Where("mark = 0 AND homework_publish_id = ?").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

//用户是否提交过作业
func HasUserSubmitHomework(userId uint, publishId uint) (bool, error) {
	var submit HomeworkSubmit
	err := db.Where("user_id = ? AND homework_publish_id = ?", userId, publishId).
		First(&submit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}
func CalculateUserHomeworkScore(submit HomeworkSubmit, tx *gorm.DB) (err error) {
	var totalScore uint
	err = tx.Where("user_id = ? AND homework_publish_id = ?",
		submit.UserId, submit.HomeworkPublishId).First(&submit).Error
	if err != nil {
		return err
	}
	if err = tx.Model(&HomeworkSubmitItem{}).Select("SUM(score)").
		Where("homework_submit_id = ?", submit.ID).Row().Scan(&totalScore); err != nil {
		return err
	}
	if err = tx.Model(&submit).Update("total_score", totalScore).Error; err != nil {
		return err
	}
	return nil

}

//创建用户提交的作业记录
func AddHomeworkSubmitWithItems(submit HomeworkSubmit) (uint, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&submit).Error; err != nil { //创建时会关联创建题目记录
		tx.Rollback()
		return 0, err
	}
	if err := CalculateUserHomeworkScore(submit, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return submit.ID, nil
}
func UpdateHomeworkSubmitWithItems(submit HomeworkSubmit) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&submit).Update(submit).Error; err != nil {
			return err
		}
		if err := CalculateUserHomeworkScore(submit, tx); err != nil {
			return err
		}
		return nil
	})
	return err

}
