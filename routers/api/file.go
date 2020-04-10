package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"learning/conf"
	"learning/service"
	"learning/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type FileForm struct {
	CourseId uint `form:"courseId" binding:"required" `
	ParentId uint `form:"parentId"  `
}

func GetChildFile(c *gin.Context) {
	var form FileForm
	if err := c.ShouldBindQuery(&form); err != nil {
		c.String(http.StatusBadRequest, "")
	} else {
		s := service.FileService{
			ParentId: form.ParentId,
			CourseId: form.CourseId,
		}
		if files, err := s.GetChildFile(); err == nil {
			c.JSON(http.StatusOK, files)
		} else {
			c.String(http.StatusInternalServerError, "")
		}
	}

}
func CreateFile(c *gin.Context) {
	var form FileForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
	} else {
		if claims, ok := c.Get("claims"); ok {
			multiForm, err := c.MultipartForm()
			if err != nil {
				c.String(http.StatusBadRequest, "")
				log.Println("未获取到文件")
				return
			}
			files := multiForm.File["file[]"]
			for _, file := range files {
				u1 := uuid.Must(uuid.NewV4(), nil).String()
				str := strings.Split(file.Filename, ".")
				suffix := ""
				if len(str) > 1 {
					suffix = "." + str[len(str)-1]
				}
				s := service.FileService{
					UserId:        claims.(*utils.Claims).Id,
					ParentId:      form.ParentId,
					CourseId:      form.CourseId,
					Filename:      file.Filename,
					Size:          uint(file.Size),
					LocalFilename: u1 + suffix,
				}
				if s.ParentId == 0 {
					s.LocalFilepath = ""
				} else {
					if f, err := s.GetParentFileByParentId(); err == nil {
						s.LocalFilepath = f.LocalFilepath + "/" + strconv.Itoa(int(f.ID))
					}
				}
				if _, err := s.AddFile(); err != nil {
					log.Println("插入数据库失败", err)
				} else {
					var directory string = conf.AppConfig.Path.File + s.LocalFilepath
					err = os.MkdirAll(directory, os.ModePerm)
					if err != nil {
						log.Println("创建目录失败", err)
					}
					localFilepath := directory + "/" + u1 + suffix
					if err := c.SaveUploadedFile(file, localFilepath); err == nil {
						log.Println("保存文件成功")
					} else {
						log.Println("保存文件失败", err)
					}
				}

			}
			c.JSON(http.StatusCreated, gin.H{})
		}
	}
}

type CreateFolderForm struct {
	CourseId   uint   `form:"courseId" binding:"required" `
	FolderName string `form:"folderName" binding:"required" `
	ParentId   int    `form:"parentId"`
}

func CreateFolder(c *gin.Context) {
	var form CreateFolderForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Println(err)
	} else {
		if claims, ok := c.Get("claims"); ok {
			s := service.FileService{
				UserId:   claims.(*utils.Claims).Id,
				ParentId: uint(form.ParentId),
				CourseId: form.CourseId,
				Filename: form.FolderName,
			}
			if s.ParentId == 0 {
				s.LocalFilepath = ""
			} else {
				if f, err := s.GetParentFileByParentId(); err == nil {
					s.LocalFilepath = f.LocalFilepath + "/" + strconv.Itoa(int(f.ID))
				}
			}
			if id, err := s.AddFile(); err == nil {
				var directory string = conf.AppConfig.Path.File +
					s.LocalFilepath + "/" + strconv.Itoa(int(id))
				err = os.MkdirAll(directory, os.ModePerm)
				if err != nil {
					log.Println("创建目录失败", err)
				}
				c.String(http.StatusCreated, "")
				return
			}
		}
		c.String(http.StatusInternalServerError, "")
	}
}
func DownloadFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.FileService{
		Id: uint(id),
	}
	file, err := s.GetFileById()
	if err != nil || file == nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	filename := url.QueryEscape(file.Filename)
	filepath := conf.AppConfig.Path.File + file.LocalFilepath + "/" + file.LocalFilename
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(filepath)
}
func DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "")
		return
	}
	s := service.FileService{
		Id: uint(id),
	}
	file, _ := s.GetFileById()
	if err := s.DeleteFileById(); err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		if file != nil {
			if file.LocalFilename == "" { //如果是目录
				err := os.RemoveAll(conf.AppConfig.Path.File +
					file.LocalFilepath + "/" + strconv.Itoa(int(file.ID)))
				if err != nil {
					log.Println("删除文件夹失败", err)
				}
			} else {
				err := os.Remove(conf.AppConfig.Path.File + file.LocalFilepath + "/" + file.LocalFilename)
				if err != nil {
					log.Println("删除文件失败", err)
				}
			}
		}
		c.String(http.StatusOK, "")
	}
}
