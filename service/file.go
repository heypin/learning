package service

import (
	"learning/models"
)

type FileService struct {
	Id            uint
	UserId        uint
	ParentId      uint
	CourseId      uint
	Filename      string
	LocalFilename string
	LocalFilepath string
	Size          uint
}

func (s *FileService) AddFile() (id uint, err error) {
	f := models.File{
		UserId:        s.UserId,
		ParentId:      s.ParentId,
		CourseId:      s.CourseId,
		Filename:      s.Filename,
		LocalFilepath: s.LocalFilepath,
		LocalFilename: s.LocalFilename,
		Size:          s.Size,
	}
	if id, err = models.AddFile(f); err != nil {
		return 0, err
	}
	return id, nil
}
func (s *FileService) GetChildFile() (f []*models.File, err error) {
	if f, err = models.GetChildFileByCourseId(s.CourseId, s.ParentId); err != nil {
		return nil, err
	}
	return f, nil
}
func (s *FileService) GetFileById() (*models.File, error) {
	f, err := models.GetFileById(s.Id)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (s *FileService) GetParentFileByParentId() (*models.File, error) {
	f, err := models.GetFileById(s.ParentId)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (s *FileService) DeleteFileById() (err error) {
	err = models.DeleteFileById(s.Id)
	return err
}
