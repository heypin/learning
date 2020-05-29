package service

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"learning/cache"
	"learning/models"
	"strconv"
	"time"
)

type NotifyService struct {
	Id       uint
	CourseId uint
	Title    string
	Content  string
}

func deleteCourseNotifyCache(courseId uint) {
	//log.Println("删除缓存")
	cache.RedisClient.Del(cache.NotifyPrefix + "." + strconv.Itoa(int(courseId)))
}
func (s *NotifyService) AddNotify() (id uint, err error) {
	n := models.Notify{
		CourseId: s.CourseId,
		Title:    s.Title,
		Content:  s.Content,
	}
	id, err = models.AddNotify(n)
	if err == nil {
		deleteCourseNotifyCache(s.CourseId)
	}
	return
}
func (s *NotifyService) GetNotifyByCourseId() (n []*models.Notify, err error) {
	key := cache.NotifyPrefix + "." + strconv.Itoa(int(s.CourseId))
	if value, err := cache.RedisClient.Get(key).Result(); err == nil {
		var cacheNotifies []*models.Notify
		if err := json.Unmarshal([]byte(value), &cacheNotifies); err == nil {
			//log.Println("使用缓存")
			return cacheNotifies, nil
		}
	}
	n, err = models.GetNotifyByCourseId(s.CourseId)
	if err == nil {
		data, err := json.Marshal(&n)
		if err == nil {
			cache.RedisClient.Set(key, string(data), time.Minute*10)
		}
	}
	return
}
func (s *NotifyService) UpdateNotifyById() (err error) {
	n := models.Notify{
		Model:   gorm.Model{ID: s.Id},
		Title:   s.Title,
		Content: s.Content,
	}
	err = models.UpdateNotifyById(n)
	if notify, _ := models.GetNotifyById(s.Id); notify != nil {
		deleteCourseNotifyCache(notify.CourseId)
	}
	return
}
func (s *NotifyService) DeleteNotifyById() (err error) {
	notify, _ := models.GetNotifyById(s.Id) //在被删除前先获取通知以便获取课程ID
	err = models.DeleteNotifyById(s.Id)
	if notify != nil {
		deleteCourseNotifyCache(notify.CourseId)
	}
	return
}
