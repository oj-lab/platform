package model

import "github.com/google/uuid"

type SubmissionStatus string

const (
	SubmissionStatusPending  SubmissionStatus = "pending"
	SubmissionStatusWaiting  SubmissionStatus = "waiting"
	SubmissionStatusRunning  SubmissionStatus = "running"
	SubmissionStatusFinished SubmissionStatus = "finished"
)

type SubmissionLanguage string

const (
	SubmissionLanguageCpp    SubmissionLanguage = "Cpp"
	SubmissionLanguageRust   SubmissionLanguage = "Rust"
	SubmissionLanguagePython SubmissionLanguage = "Python"
)

// Using relationship according to https://gorm.io/docs/belongs_to.html
type JudgeTaskSubmission struct {
	MetaFields
	UID         uuid.UUID          `gorm:"primaryKey" json:"uid"`
	UserAccount string             `gorm:"not null" json:"userAccount"`
	User        User               `json:"user"`
	ProblemSlug string             `gorm:"not null" json:"problemSlug"`
	Problem     Problem            `json:"problem"`
	Code        string             `gorm:"not null" json:"code"`
	Language    SubmissionLanguage `gorm:"not null" json:"language"`
	Status      SubmissionStatus   `gorm:"default:pending" json:"status"`
	VerdictJson string             `json:"verdictJson"`
}

type JudgeTaskSubmissionSortByColumn string

const (
	JudgeTaskSubmissionSortByColumnCreateAt JudgeTaskSubmissionSortByColumn = "create_at"
	JudgeTaskSubmissionSortByColumnUpdateAt JudgeTaskSubmissionSortByColumn = "update_at"
)

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
