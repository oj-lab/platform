package judge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

func PickJudgeTask(ctx context.Context, consumer string) (*judge_model.JudgeTask, error) {
	task, err := getTaskFromStream(ctx, consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to get task from stream: %w", err)
	}

	db := gorm_agent.GetDefaultDB()
	err = judge_model.UpdateSubmission(db, judge_model.JudgeTaskSubmission{
		UID:    uuid.MustParse(task.SubmissionUID),
		Status: judge_model.SubmissionStatusRunning,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update submission status: %w", err)
	}

	return task, nil
}

func ReportJudgeTaskResult(
	ctx context.Context,
	consumer string, streamID string, verdictJson string,
) error {
	db := gorm_agent.GetDefaultDB()

	mainVerdict, err := parseVerdictJson(verdictJson)
	if err != nil {
		return err
	}
	err = judge_model.UpdateSubmission(db, judge_model.JudgeTaskSubmission{
		RedisStreamID: streamID,
		Status:        judge_model.SubmissionStatusFinished,
		VerdictJson:   verdictJson,
		MainResult:    mainVerdict,
	})

	if err != nil {
		return err
	}

	err = ackTaskFromStream(ctx, consumer, streamID)
	if err != nil {
		return err
	}

	return nil
}

type VerdictJson struct {
	Verdict           string  `json:"verdict"`
	TimeUsage         Time    `json:"time_usage"`
	MemoryUsageBytes  float64 `json:"memory_usage_bytes"`
	ExitStatus        int     `json:"exit_status"`
	CheckerExitStatus int     `json:"checker_exit_status"`
}

type Time struct {
	Secs  float64 `json:"secs"`
	Nanos float64 `json:"nanos"`
}

func parseVerdictJson(verdictString string) (judge_model.JudgeVerdict, error) {
	var tests []VerdictJson
	err := json.Unmarshal([]byte(verdictString), &tests)
	if err != nil {
		return "", err
	}

	var (
		totolTestPoint = len(tests)
		Priority       = 6
		AvgMiles       = 0.0
		MaxMiles       = 0.0
		AvgMemoryBytes = 0.0
		MaxMemoryBytes = 0.0
	)
	verdictPriorityMap := map[judge_model.JudgeVerdict]int{
		judge_model.JudgeVerdictCompileError:        0,
		judge_model.JudgeVerdictRuntimeError:        1,
		judge_model.JudgeVerdictTimeLimitExceeded:   2,
		judge_model.JudgeVerdictMemoryLimitExceeded: 3,
		judge_model.JudgeVerdictSystemError:         4,
		judge_model.JudgeVerdictWrongAnswer:         5,
		judge_model.JudgeVerdictAccepted:            6,
	}
	priorityVerdictMap := map[int]judge_model.JudgeVerdict{
		0: judge_model.JudgeVerdictCompileError,
		1: judge_model.JudgeVerdictRuntimeError,
		2: judge_model.JudgeVerdictTimeLimitExceeded,
		3: judge_model.JudgeVerdictMemoryLimitExceeded,
		4: judge_model.JudgeVerdictSystemError,
		5: judge_model.JudgeVerdictWrongAnswer,
		6: judge_model.JudgeVerdictAccepted,
	}

	for _, test := range tests {
		tempMiles := test.TimeUsage.Secs*1000 + test.TimeUsage.Nanos/1000000
		Priority = min(Priority, verdictPriorityMap[judge_model.JudgeVerdict(test.Verdict)])
		AvgMiles += tempMiles
		MaxMiles = max(MaxMiles, tempMiles)
		AvgMemoryBytes += test.MemoryUsageBytes
		MaxMemoryBytes = max(MaxMemoryBytes, test.MemoryUsageBytes)
	}

	AvgMiles /= float64(totolTestPoint)
	AvgMemoryBytes /= float64(totolTestPoint)
	finalVerdict := priorityVerdictMap[Priority]

	log.AppLogger().Debugln(totolTestPoint, finalVerdict, AvgMiles, MaxMiles, AvgMemoryBytes, MaxMemoryBytes)

	// models.JudgeResult{
	// 	MainVerdict:    finalVerdict,
	// 	TestPointCount: uint64(totolTestPoint),
	// 	MaxTimeMs:      uint64(MaxMiles),
	// 	AverageTimeMs:  uint64(AvgMiles),
	// 	maxTimeMs:      uint64(MaxTime),
	// 	AverageMemory:  uint64(AvgMemoryBytes),
	// 	MaxMemory:      uint64(MaxMemoryBytes),
	// }
	return finalVerdict, nil
}
