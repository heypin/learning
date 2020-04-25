package models

import "github.com/jinzhu/gorm"

type HomeworkSubmit struct {
	gorm.Model
	UserId            uint                  `json:"userId"`
	HomeworkPublishId uint                  `json:"homeworkPublishId"`
	TotalScore        uint                  `json:"totalScore"`
	Mark              *uint                 `json:"mark"`
	SubmitItems       []*HomeworkSubmitItem `json:"submitItems" gorm:"foreignKey:homework_submit_id;association_foreignKey:id;"`
	User              User                  `json:"user" gorm:"foreignKey:id;association_foreignKey:user_id;"`
}

func GetHomeworkUserSubmitWithItems(userId uint, publishId uint) (*HomeworkSubmit, error) {
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
func UpdateHomeworkSubmitById(submit HomeworkSubmit) error {
	if err := db.Model(&submit).Update(submit).Error; err != nil {
		return err
	}
	return nil
}
func GetHomeworkSubmitById(id uint) (*HomeworkSubmit, error) {
	var submit HomeworkSubmit
	err := db.Where("id = ?", id).First(&submit).Error
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
	err = db.Where("homework_publish_id = ?", publishId).
		Preload("User").Find(&submits).Error
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
	err = db.Model(&HomeworkSubmit{}).Where("mark = 0 AND homework_publish_id = ?", publishId).
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
	if err = tx.Model(&HomeworkSubmitItem{}).Select("SUM(score)").
		Where("homework_submit_id = ?", submit.ID).Row().Scan(&totalScore); err != nil {
		return err
	}
	if err = tx.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Model(&submit).Update("total_score", totalScore).Error; err != nil {
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
		for _, item := range submit.SubmitItems {
			item.HomeworkSubmitId = submit.ID
			if item.ID == 0 { //有ID更新它，无主键创建它
				if err := tx.Create(item).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(item).Update(*item).Error; err != nil {
					return err
				}
			}
		}

		submit.SubmitItems = nil //子项已在上一步更新，gorm关联更新会更新关联结构体零值,这里置为空,或者使用
		//Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).
		if err := tx.Model(&submit).Updates(submit).Error; err != nil {
			return err
		}
		if err := CalculateUserHomeworkScore(submit, tx); err != nil {
			return err
		}
		return nil
	})
	return err

}
