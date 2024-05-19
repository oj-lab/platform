package judge

import (
	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/models"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
)

type SubmissionStatus string

const (
	SubmissionStatusPending  SubmissionStatus = "pending"
	SubmissionStatusWaiting  SubmissionStatus = "waiting"
	SubmissionStatusRunning  SubmissionStatus = "running"
	SubmissionStatusFinished SubmissionStatus = "finished"
)

type SubmissionLanguage string

func (sl SubmissionLanguage) String() string {
	return string(sl)
}

const (
	SubmissionLanguageCpp    SubmissionLanguage = "Cpp"
	SubmissionLanguageRust   SubmissionLanguage = "Rust"
	SubmissionLanguagePython SubmissionLanguage = "Python"
)

// Using relationship according to https://gorm.io/docs/belongs_to.html
type JudgeTaskSubmission struct {
	models.MetaFields
	UID           uuid.UUID             `gorm:"primaryKey" json:"UID"`
	RedisStreamID string                `json:"redisStreamID"`
	UserAccount   string                `gorm:"not null" json:"userAccount"`
	User          user_model.User       `json:"user"`
	ProblemSlug   string                `gorm:"not null" json:"problemSlug"`
	Problem       problem_model.Problem `json:"problem"`
	Code          string                `gorm:"not null" json:"code"`
	Language      SubmissionLanguage    `gorm:"not null" json:"language"`
	Status        SubmissionStatus      `gorm:"default:pending" json:"status"`
	VerdictJson   string                `json:"verdictJson"`
	MainResult    JudgeVerdict          `json:"mainResult"`
}

func NewSubmission(
	userAccount string,
	problemSlug string,
	code string,
	language SubmissionLanguage,
) JudgeTaskSubmission {
	return JudgeTaskSubmission{
		UserAccount: userAccount,
		ProblemSlug: problemSlug,
		Code:        code,
		Language:    language,
		Status:      SubmissionStatusPending,
	}
}

func (s *JudgeTaskSubmission) ToJudgeTask() JudgeTask {
	return JudgeTask{
		SubmissionUID: s.UID.String(),
		ProblemSlug:   s.ProblemSlug,
		Code:          s.Code,
		Language:      s.Language.String(),
	}
}
