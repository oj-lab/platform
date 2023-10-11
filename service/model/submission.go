package model

type SubmissionStatus string

const (
	SubmissionStatusPending SubmissionStatus = "pending"
	SubmissionStatusRunning SubmissionStatus = "running"
	SubmissionStatusDone    SubmissionStatus = "done"
)

// Using relationship according to https://gorm.io/docs/belongs_to.html
type Submission struct {
	MetaFields
	UID         string           `gorm:"primaryKey" json:"uid"`
	UserAccount string           `gorm:"not null" json:"userAccount"`
	User        User             `json:"user"`
	ProblemSlug string           `gorm:"not null" json:"problemSlug"`
	Problem     Problem          `json:"problem"`
	Status      SubmissionStatus `gorm:"not null"`
}
