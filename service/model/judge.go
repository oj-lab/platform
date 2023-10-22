package model

import (
	"strings"

	"github.com/google/uuid"
)

type JudgerState string

const (
	JudgerStateIdle    JudgerState = "idle"
	JudgerStateBusy    JudgerState = "busy"
	JudgerStateOffline JudgerState = "offline"
)

type Judger struct {
	MetaFields
	Host  string      `gorm:"primaryKey" json:"host"`
	State JudgerState `gorm:"default:offline" json:"status"`
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

func (js JudgerState) CanUpdate(nextStatus JudgerState) bool {
	switch js {
	case JudgerStateOffline:
		return nextStatus == JudgerStateIdle
	case JudgerStateIdle:
		return nextStatus == JudgerStateBusy || nextStatus == JudgerStateOffline
	case JudgerStateBusy:
		return nextStatus == JudgerStateIdle || nextStatus == JudgerStateOffline
	default:
		return false
	}
}

func StringToJudgerState(state string) JudgerState {
	state = strings.ToLower(state)
	switch state {
	case "idle":
		return JudgerStateIdle
	case "busy":
		return JudgerStateBusy
	case "offline":
		return JudgerStateOffline
	default:
		return JudgerStateOffline
	}
}
