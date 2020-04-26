package api

import (
	"github.com/gin-gonic/gin"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"strconv"
)

func GetCommentByCourseId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.CommentService{
		CourseId: uint(courseId),
	}
	if comments, err := s.GetCommentByCourseId(); err == nil {
		c.JSON(http.StatusOK, comments)
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
func GetCommentByUserId(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.CommentService{
			UserId:   claims.(*utils.Claims).Id,
			CourseId: uint(courseId),
		}
		if comments, err := s.GetCommentByUserId(); err == nil {
			c.JSON(http.StatusOK, comments)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
func GetCommentReplyToUser(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Query("courseId"))
	if err != nil || courseId <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.CommentService{
			UserId:   claims.(*utils.Claims).Id,
			CourseId: uint(courseId),
		}
		if comments, err := s.GetCommentReplyToUser(); err == nil {
			c.JSON(http.StatusOK, comments)
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}

type CreateCommentForm struct {
	CourseId    uint   `json:"courseId" binding:"required"`
	ParentId    uint   `json:"parentId" `
	ReplyId     uint   `json:"replyId" `
	ReplyUserId uint   `json:"replyUserId"`
	Content     string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	var form CreateCommentForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
		return
	}
	if claims, ok := c.Get("claims"); ok {
		s := service.CommentService{
			CourseId:    form.CourseId,
			ParentId:    form.ParentId,
			ReplyId:     form.ReplyId,
			ReplyUserId: form.ReplyUserId,
			UserId:      claims.(*utils.Claims).Id,
			Content:     form.Content,
		}
		if _, err := s.AddComment(); err == nil {
			c.String(http.StatusCreated, "")
			return
		}
	}
	c.String(http.StatusInternalServerError, "")
}
func DeleteCommentById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.CommentService{
		Id: uint(id),
	}
	if err = s.DeleteCommentById(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
