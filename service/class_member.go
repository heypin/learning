package service

import (
	"errors"
	"learning/models"
)

type ClassMemberService struct {
	ClassId uint
	UserId  uint
}

func (s *ClassMemberService) GetUsersByClassId() (users []*models.User, err error) {
	users, err = models.GetUsersByClassId(s.ClassId)
	return
}
func (s *ClassMemberService) GetClassesByUserId() (classes []*models.Class, err error) {
	classes, err = models.GetClassesByUserId(s.UserId)
	return
}
func (s *ClassMemberService) JoinClassByClassCode(code string) error {
	if class, err := models.GetClassByClassCode(code); err == nil && class != nil {
		if ok, err := models.HasJoinClass(s.UserId, class.ID); !ok && err == nil {
			if err := models.AddClassMember(s.UserId, s.ClassId); err == nil {
				return nil
			} else {
				return err
			}
		} else if ok {
			return nil
		} else {
			return err
		}
	} else if err == nil && class == nil {
		return errors.New("该班级不存在")
	} else {
		return err
	}
}
func (s *ClassMemberService) DeleteClassMember() (err error) {
	err = models.DeleteClassMember(s.UserId, s.ClassId)
	return
}
