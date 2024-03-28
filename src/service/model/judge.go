package model

import (
	"strings"
)

// Should contains a priority definition
// Ex. CompileError > RuntimeError > TimeLimitExceeded > MemoryLimitExceeded > SystemError > WrongAnswer > Accepted
type JudgeVerdict string

const (
	JudgeVerdictCompileError        JudgeVerdict = "CompileError" // Only for main verdict
	JudgeVerdictRuntimeError        JudgeVerdict = "RuntimeError"
	JudgeVerdictTimeLimitExceeded   JudgeVerdict = "TimeLimitExceeded"
	JudgeVerdictMemoryLimitExceeded JudgeVerdict = "MemoryLimitExceeded"
	JudgeVerdictSystemError         JudgeVerdict = "SystemError" // Some runtime unknown error ?
	JudgeVerdictWrongAnswer         JudgeVerdict = "WrongAnswer"
	JudgeVerdictAccepted            JudgeVerdict = "Accepted"
	JudgeVerdictCancelled           JudgeVerdict = "cancelled" // Judge will be cancelled if some point results in Runtime error, Time limit exceeded, Memory limit exceeded
)

type JudgeResult struct {
	MainVerdict    JudgeVerdict         `json:"verdict"`        // A merge of all TestPoints' verdict, according to the pirority
	Detail         string               `json:"detail"`         // A brief description of the result
	TestPointCount uint64               `json:"testPointCount"` // Won't be stored in database
	TestPointMap   map[string]TestPoint `json:"testPoints"`     // Won't be stored in database
	TestPointsJson string               `json:"-"`              // Used to store TestPoints in database
	AverageTimeMs  uint64               `json:"averageTimeMs"`  // Won't be stored in database
	MaxTimeMs      uint64               `json:"maxTimeMs"`      // Won't be stored in database
	AverageMemory  uint64               `json:"averageMemory"`  // Won't be stored in database
	MaxMemory      uint64               `json:"maxMemory"`      // Won't be stored in database
}

type TestPoint struct {
	Index           string       `json:"index"` // The name of *.in/ans file
	Verdict         JudgeVerdict `json:"verdict"`
	Diff            *ResultDiff  `json:"diff"` // Required if verdict is wrong_answer
	TimeUsageMs     uint64       `json:"timeUsageMs"`
	MemoryUsageByte uint64       `json:"memoryUsageByte"`
}

type ResultDiff struct {
	Expected string `json:"expected"`
	Received string `json:"received"`
}

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
