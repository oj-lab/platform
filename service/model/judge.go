package model

import "github.com/google/uuid"

type Judger struct {
	MetaFields
	Host string `gorm:"primaryKey" json:"host"`
}

type JudgeTask struct {
	UID         uuid.UUID `json:"uid"`
	ProblemSlug string    `json:"problemSlug"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Judger      Judger    `json:"judger"`
}

func NewJudgeTask(problemSlug, code, language string) *JudgeTask {
	return &JudgeTask{
		UID:         uuid.New(),
		ProblemSlug: problemSlug,
		Code:        code,
		Language:    language,
	}
}
