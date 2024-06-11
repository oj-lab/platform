package judge_model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/models"
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
	JudgeVerdictCancelled           JudgeVerdict = "Cancelled" // Judge will be cancelled if some point results in Runtime error, Time limit exceeded, Memory limit exceeded
)

func (jv JudgeVerdict) IsValid() bool {
	switch jv {
	case JudgeVerdictCompileError,
		JudgeVerdictRuntimeError,
		JudgeVerdictTimeLimitExceeded,
		JudgeVerdictMemoryLimitExceeded,
		JudgeVerdictSystemError,
		JudgeVerdictWrongAnswer,
		JudgeVerdictAccepted,
		JudgeVerdictCancelled:
		return true
	}
	return false
}

var ErrInvalidJudgeVerdict = fmt.Errorf("invalid JudgeVerdict")

type JudgeResult struct {
	models.MetaFields
	UID             uuid.UUID    `json:"UID" gorm:"primaryKey"`
	JudgeUID        uuid.UUID    `json:"judgeUID"`
	Verdict         JudgeVerdict `json:"verdict"`
	TimeUsageMS     uint         `json:"timeUsageMS"`
	MemoryUsageByte uint         `json:"memoryUsageByte"`
	Output          string       `json:"output"`
	ExpectedOutput  string       `json:"expectedOutput"`
	SystemOutput    string       `json:"systemOutput"`
}
