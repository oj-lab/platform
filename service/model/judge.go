package model

import "github.com/google/uuid"

type JudgeTask struct {
	UID         string `json:"uid"`
	ProblemSlug string `json:"problemSlug"`
	Src         string `json:"src"`
	SrcLanguage string `json:"srcLanguage"`
}

func NewJudgeTask(problemSlug, src, srcLanguage string) *JudgeTask {
	return &JudgeTask{
		UID:         uuid.NewString(),
		ProblemSlug: problemSlug,
		Src:         src,
		SrcLanguage: srcLanguage,
	}
}
