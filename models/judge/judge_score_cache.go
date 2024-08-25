package judge_model

import (
	"time"

	"github.com/oj-lab/oj-lab-platform/models"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
)

// user contest problem summary score info
type ScoreCache struct {
	models.MetaFields
	UserAccount     string                `json:"userAccount" gorm:"primaryKey"`
	User            user_model.User       `json:"user"`
	ProblemSlug     string                `json:"problemSlug" gorm:"primaryKey"`
	Problem         problem_model.Problem `json:"problem"`
	SubmissionCount int64                 `json:"submissionCount" gorm:"default:1"` // judge create time < solvetime will be count
	IsCorrect       bool                  `json:"isCorrect" gorm:"default:false"`
	SolveTime       *time.Time            `json:"SolveAt"  gorm:"default:null"` // ac time < solveTime, update submissionCount
}

func NewScoreCache(userAccount string, problemSlug string) ScoreCache {
	return ScoreCache{
		UserAccount:     userAccount,
		ProblemSlug:     problemSlug,
		SubmissionCount: 1,
		IsCorrect:       false,
	}
}
