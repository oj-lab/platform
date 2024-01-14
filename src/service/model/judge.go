package model

import (
	"strings"
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
	SubmissionUID string  `json:"submissionUID"`
	ProblemSlug   string  `json:"problemSlug"`
	Code          string  `json:"code"`
	Language      string  `json:"language"`
	RedisStreamID *string `json:"redisStreamID"`
}

func (jt *JudgeTask) ToStringMap() map[string]interface{} {
	return map[string]interface{}{
		"submission_uid": jt.SubmissionUID,
		"problem_slug":   jt.ProblemSlug,
		"code":           jt.Code,
		"language":       jt.Language,
	}
}

func JudgeTaskFromMap(m map[string]interface{}) *JudgeTask {
	return &JudgeTask{
		SubmissionUID: m["submission_uid"].(string),
		ProblemSlug:   m["problem_slug"].(string),
		Code:          m["code"].(string),
		Language:      m["language"].(string),
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
