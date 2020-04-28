package service

import (
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
	if class, err := models.GetClassByClassCode(code); err == nil {
		if ok, err := models.ExistClassMemberRecord(s.UserId, class.ID); ok {
			return nil
		} else if err == nil {
			if err := models.AddClassMember(s.UserId, class.ID); err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
func (s *ClassMemberService) DeleteClassMember() (err error) {
	err = models.DeleteClassMember(s.UserId, s.ClassId)
	return
}
