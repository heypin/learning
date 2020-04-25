package models

type HomeworkLibItemOption struct {
	HomeworkLibItemId uint   `json:"homeworkLibItemId"`
	Sequence          string `json:"sequence"`
	Content           string `json:"content"`
}
