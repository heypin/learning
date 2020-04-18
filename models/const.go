package models

const (
	ROLE_ADMIN = iota
	ROLE_USER
)
const (
	SEX_MALE   = 1
	SEX_FEMALE = 2
)
const (
	Resubmit_Deny  = 0
	Resubmit_Allow = 1
)
const (
	Subject_Single    = "单选题"
	Subject_Multiple  = "多选题"
	Subject_Judgement = "判断题"
	Subject_Blank     = "填空题"
	Subject_Short     = "简答题"
	Subject_Program   = "编程题"
)
