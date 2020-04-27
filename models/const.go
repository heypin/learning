package models

import "errors"

var ErrRecordExist = errors.New("记录已存在")

const (
	RoleAdmin = iota
	RoleUser
)

//const (
//SEX_MALE   = 1
//SEX_FEMALE = 2
//)
const (
	//Resubmit_Deny  = 0
	ResubmitAllow = 1
)
const (
	SubjectSingle    = "单选题"
	SubjectMultiple  = "多选题"
	SubjectJudgement = "判断题"
	SubjectBlank     = "填空题"
	SubjectShort     = "简答题"
	SubjectProgram   = "编程题"
)
