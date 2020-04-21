package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"learning/models"
)

type HomeworkSubmitService struct {
	Id                uint
	UserId            uint
	HomeworkPublishId uint
	TotalScore        uint
	Mark              *uint
	SubmitItems       []*models.HomeworkSubmitItem
}

//获取用户提交的作业及其答案
func (s *HomeworkSubmitService) GetHomeworkSubmitWithItems() (*models.HomeworkSubmit, error) {
	submit, err := models.GetHomeworkSubmitWithItems(s.UserId, s.HomeworkPublishId)
	return submit, err
}

//创建提交作业记录，如果已创建，则表示进行更新
func (s *HomeworkSubmitService) SubmitHomeworkWithItems() (uint, error) {
	if ok, err := models.HasUserSubmitHomework(s.UserId, s.HomeworkPublishId); ok {
		if allow, _ := models.IsAllowResubmitHomework(s.HomeworkPublishId); allow {
			if err := s.UpdateSubmitHomeworkWithItems(); err != nil {
				return 0, err
			} else {
				return 0, nil
			}
		} else {
			return 0, errors.New("不允许重复提交")
		}

	} else if !ok && err != nil {
		return 0, err
	} else {
		submit := models.HomeworkSubmit{
			UserId:            s.UserId,
			HomeworkPublishId: s.HomeworkPublishId,
			TotalScore:        s.TotalScore,
			Mark:              s.Mark,
			SubmitItems:       s.SubmitItems,
		}
		id, err := models.AddHomeworkSubmitWithItems(submit)
		return id, err
	}
}

//更新提交
func (s *HomeworkSubmitService) UpdateSubmitHomeworkWithItems() error {
	submit := models.HomeworkSubmit{
		Model:             gorm.Model{ID: s.Id},
		UserId:            s.UserId,
		HomeworkPublishId: s.HomeworkPublishId,
		TotalScore:        s.TotalScore,
		Mark:              s.Mark,
		SubmitItems:       s.SubmitItems,
	}
	err := models.UpdateHomeworkSubmitWithItems(submit)
	return err
}

//获取发布的作业的所有提交记录
func (s *HomeworkSubmitService) GetHomeworkSubmitsByPublishId() (submits []*models.HomeworkSubmit, err error) {
	submits, err = models.GetHomeworkSubmitsByPublishId(s.HomeworkPublishId)
	return
}
