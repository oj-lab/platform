package model

type SubmissionStatus string

const (
	SubmissionStatusPending SubmissionStatus = "pending"
	SubmissionStatusRunning SubmissionStatus = "running"
	SubmissionStatusDone    SubmissionStatus = "done"
)

// Using relationship according to https://gorm.io/docs/belongs_to.html
type JudgeTaskSubmission struct {
	MetaFields
	UID         string           `gorm:"primaryKey" json:"uid"`
	UserAccount string           `gorm:"not null" json:"userAccount"`
	User        User             `json:"user"`
	ProblemSlug string           `gorm:"not null" json:"problemSlug"`
	Problem     Problem          `json:"problem"`
	Code        string           `gorm:"not null" json:"code"`
	Language    string           `gorm:"not null" json:"language"`
	Status      SubmissionStatus `gorm:"not null" json:"status"`
}

func (jts JudgeTaskSubmission) GetJudgeTask() JudgeTask {
	return JudgeTask{
		UID:         jts.UID,
		ProblemSlug: jts.ProblemSlug,
		Src:         jts.Code,
		SrcLanguage: jts.Language,
	}
}
