package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"learning/models"
	"time"
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
func (s *HomeworkSubmitService) GetHomeworkUserSubmitWithItems() (*models.HomeworkSubmit, error) {
	submit, err := models.GetHomeworkUserSubmitWithItems(s.UserId, s.HomeworkPublishId)
	return submit, err
}

//创建提交作业记录，如果已创建，则表示进行更新
func (s *HomeworkSubmitService) SubmitHomeworkWithItems() (uint, error) {
	publish, err := models.GetHomeworkPublishById(s.HomeworkPublishId)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	if now.Before(publish.BeginTime) || now.After(publish.EndTime) {
		return 0, errors.New("不在允许提交的时间范围内")
	}
	if ok, err := models.HasUserSubmitHomework(s.UserId, s.HomeworkPublishId); ok {
		if submit, _ := s.GetHomeworkSubmitById(); *submit.Mark == 1 {
			return 0, errors.New("已批阅，不能再提交")
		}
		if *publish.Resubmit == models.ResubmitAllow {
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
		Model: gorm.Model{ID: s.Id},
		//UserId:            s.UserId,
		//HomeworkPublishId: s.HomeworkPublishId,//更新时不需要
		//TotalScore:        s.TotalScore,//总分由数据库计算
		Mark:        s.Mark,
		SubmitItems: s.SubmitItems,
	}
	err := models.UpdateHomeworkSubmitWithItems(submit)
	return err
}

func (s *HomeworkSubmitService) GetHomeworkSubmitById() (submit *models.HomeworkSubmit, err error) {
	return models.GetHomeworkSubmitById(s.Id)
}

//获取发布的作业的所有提交记录
func (s *HomeworkSubmitService) GetHomeworkSubmitsByPublishId() (submits []*models.HomeworkSubmit, err error) {
	submits, err = models.GetHomeworkSubmitsByPublishId(s.HomeworkPublishId)
	return
}
func (s *HomeworkSubmitService) UpdateHomeworkSubmitMarkById() (err error) {
	submit := models.HomeworkSubmit{
		Model: gorm.Model{ID: s.Id},
		Mark:  s.Mark,
	}
	err = models.UpdateHomeworkSubmitById(submit)
	return
}
