package model

import "github.com/google/uuid"

type JudgeTask struct {
	UID         string `json:"uid"`
	ProblemSlug string `json:"problemSlug"`
	// TODO: Change to name Code
	Src string `json:"src"`
	// TODO: Change to name Language
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
