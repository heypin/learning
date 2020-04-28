package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"learning/models"
	"learning/utils"
)

type UserService struct {
	Id       uint
	Email    string
	Password string
	RealName string
	Number   string
	Avatar   string
	Sex      uint
}

func (s *UserService) Auth() (uint, bool) {
	user, _ := models.GetUserByEmail(s.Email)
	if user != nil && utils.CheckPassword(user.Password, s.Password) {
		return user.ID, true
	}
	return 0, false
}
func (s *UserService) Register() (uint, error) {
	u, err := models.GetUserByEmail(s.Email)
	if err != nil {
		return 0, err
	} else if u == nil {
		user := models.User{
			Email:    s.Email,
			Password: utils.Encrypt(s.Password),
			RealName: s.RealName,
		}
		if id, err := models.AddUser(user); err == nil {
			return id, nil
		} else {
			return 0, err
		}
	} else {
		return 0, models.ErrRecordExist
	}
}

func (s *UserService) GetUserByEmail() (user *models.User, err error) {
	user, err = models.GetUserByEmail(s.Email)
	return
}
func (s *UserService) GetUserById() (user *models.User, err error) {
	user, err = models.GetUserById(s.Id)
	return
}
func (s *UserService) UpdateUserPassword(oldPassword string) error {
	result, _ := models.GetUserById(s.Id)
	if result == nil || !utils.CheckPassword(result.Password, oldPassword) {
		return errors.New("修改失败")
	}
	user := models.User{
		Model:    gorm.Model{ID: s.Id},
		Password: utils.Encrypt(s.Password),
	}
	if err := models.UpdateUserById(user); err != nil {
		return err
	}
	return nil
}
func (s *UserService) UpdateUserById() error {
	user := models.User{
		Model:    gorm.Model{ID: s.Id},
		Email:    s.Email,
		Password: utils.Encrypt(s.Password),
		RealName: s.RealName,
		Number:   s.Number,
		Sex:      s.Sex,
		Avatar:   s.Avatar,
	}
	if err := models.UpdateUserById(user); err != nil {
		return err
	}
	return nil
}
