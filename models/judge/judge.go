package judge_model

import (
	"github.com/google/uuid"
	"github.com/oj-lab/platform/models"
	problem_model "github.com/oj-lab/platform/models/problem"
	user_model "github.com/oj-lab/platform/models/user"
)

type JudgeStatus string

const (
	JudgeStatusPending  JudgeStatus = "pending"
	JudgeStatusWaiting  JudgeStatus = "waiting"
	JudgeStatusRunning  JudgeStatus = "running"
	JudgeStatusFinished JudgeStatus = "finished"
)

type ProgrammingLanguage string

func (sl ProgrammingLanguage) String() string {
	return string(sl)
}

const (
	ProgrammingLanguageCpp    ProgrammingLanguage = "Cpp"
	ProgrammingLanguageRust   ProgrammingLanguage = "Rust"
	ProgrammingLanguagePython ProgrammingLanguage = "Python"
)

func (sl ProgrammingLanguage) IsValid() bool {
	switch sl {
	case ProgrammingLanguageCpp, ProgrammingLanguageRust, ProgrammingLanguagePython:
		return true
	}
	return false
}

// Using relationship according to https://gorm.io/docs/belongs_to.html
type Judge struct {
	models.MetaFields
	UID           uuid.UUID             `json:"UID" gorm:"primaryKey"`
	RedisStreamID string                `json:"redisStreamID" gorm:"index:idx_redis_stream_id"`
	UserAccount   string                `json:"userAccount" gorm:"not null"`
	User          user_model.User       `json:"user"`
	ProblemSlug   string                `json:"problemSlug" gorm:"not null"`
	Problem       problem_model.Problem `json:"problem"`
	Code          string                `json:"code" gorm:"not null"`
	Language      ProgrammingLanguage   `json:"language" gorm:"not null"`
	Status        JudgeStatus           `json:"status" gorm:"default:pending"`
	ResultCount   uint                  `json:"resultCount"`
	Results       []JudgeResult         `json:"results" gorm:"foreignKey:JudgeUID"`
	Verdict       JudgeVerdict          `json:"verdict"`
}

func NewJudge(
	userAccount string,
	problemSlug string,
	code string,
	language ProgrammingLanguage,
) Judge {
	return Judge{
		UserAccount: userAccount,
		ProblemSlug: problemSlug,
		Code:        code,
		Language:    language,
		Status:      JudgeStatusPending,
	}
}

func (s *Judge) ToJudgeTask() JudgeTask {
	return JudgeTask{
		JudgeUID:    s.UID.String(),
		ProblemSlug: s.ProblemSlug,
		Code:        s.Code,
		Language:    s.Language.String(),
	}
}
