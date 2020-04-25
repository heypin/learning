package models

type ExamLibItemOption struct {
	ExamLibItemId uint   `json:"examLibItemId"`
	Sequence      string `json:"sequence"`
	Content       string `json:"content"`
}
