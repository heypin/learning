package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"net/http"
	"strconv"
)

func GetUsersByClassId(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ClassMemberService{
		ClassId: uint(classId),
	}
	if users, err := s.GetUsersByClassId(); err == nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}

func GetClassesByUserId(c *gin.Context) {
	if claims, ok := c.Get("claims"); ok {
		s := service.ClassMemberService{
			UserId: claims.(*utils.Claims).Id,
		}
		if classes, err := s.GetClassesByUserId(); err == nil {
			c.JSON(http.StatusOK, classes)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

func JoinClassByClassCode(c *gin.Context) {
	form := struct {
		ClassCode string `json:"classCode" binding:"required"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.ClassMemberService{
			UserId: claims.(*utils.Claims).Id,
		}
		if err := s.JoinClassByClassCode(form.ClassCode); err == nil {
			c.String(http.StatusOK, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

type ClassMemberForm struct {
	ClassId uint `json:"classId" binding:"required"`
	UserId  uint `json:"userId" binding:"userId"`
}

func DeleteClassMember(c *gin.Context) {
	classId, err := strconv.Atoi(c.Query("classId"))
	if err != nil || classId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil || userId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.ClassMemberService{
		UserId:  uint(userId),
		ClassId: uint(classId),
	}
	if err := s.DeleteClassMember(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
