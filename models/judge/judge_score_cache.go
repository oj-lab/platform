package judge_model

import (
	"time"

	"github.com/oj-lab/oj-lab-platform/models"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
)

// user contest problem summary score info
type JudgeScoreCache struct {
	models.MetaFields
	UserAccount     string                `json:"userAccount" gorm:"primaryKey"`
	User            user_model.User       `json:"user"`
	ProblemSlug     string                `json:"problemSlug" gorm:"primaryKey"`
	Problem         problem_model.Problem `json:"problem"`
	SubmissionCount int64                 `json:"submissionCount" gorm:"default:1"` // judge create time < solvetime will be count
	IsAccepted      bool                  `json:"isAccepted" gorm:"default:false"`
	SolveTime       *time.Time            `json:"solveAt"  gorm:"default:null"` // ac time < solveTime, update submissionCount
}

func NewJudgeScoreCache(userAccount string, problemSlug string) JudgeScoreCache {
	return JudgeScoreCache{
		UserAccount:     userAccount,
		ProblemSlug:     problemSlug,
		SubmissionCount: 1,
		IsAccepted:      false,
	}
}
