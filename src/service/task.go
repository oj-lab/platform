package service

import (
	"context"
	"encoding/json"
	"github.com/OJ-lab/oj-lab-services/src/core"
	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/business"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"github.com/google/uuid"
)

func PickJudgeTask(ctx context.Context, consumer string) (*model.JudgeTask, error) {
	task, err := business.GetTaskFromStream(ctx, consumer)
	if err != nil {
		return nil, err
	}

	db := gormAgent.GetDefaultDB()
	err = mapper.UpdateSubmission(db, model.JudgeTaskSubmission{
		UID:    uuid.MustParse(task.SubmissionUID),
		Status: model.SubmissionStatusRunning,
	})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func ReportJudgeTaskResult(
	ctx context.Context,
	consumer string, streamID string, verdictJson string,
) error {
	db := gormAgent.GetDefaultDB()

	mainVerdict, err := parseVerdictJson(verdictJson)
	if err != nil {
		return err
	}
	err = mapper.UpdateSubmission(db, model.JudgeTaskSubmission{
		RedisStreamID: streamID,
		Status:        model.SubmissionStatusFinished,
		VerdictJson:   verdictJson,
		MainResult:    mainVerdict,
	})

	if err != nil {
		return err
	}

	err = business.AckTaskFromStream(ctx, consumer, streamID)
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

func parseVerdictJson(verdictString string) (model.JudgeVerdict, error) {
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
	verdictPriorityMap := map[model.JudgeVerdict]int{
		model.JudgeVerdictCompileError:        0,
		model.JudgeVerdictRuntimeError:        1,
		model.JudgeVerdictTimeLimitExceeded:   2,
		model.JudgeVerdictMemoryLimitExceeded: 3,
		model.JudgeVerdictSystemError:         4,
		model.JudgeVerdictWrongAnswer:         5,
		model.JudgeVerdictAccepted:            6,
	}
	priorityVerdictMap := map[int]model.JudgeVerdict{
		0: model.JudgeVerdictCompileError,
		1: model.JudgeVerdictRuntimeError,
		2: model.JudgeVerdictTimeLimitExceeded,
		3: model.JudgeVerdictMemoryLimitExceeded,
		4: model.JudgeVerdictSystemError,
		5: model.JudgeVerdictWrongAnswer,
		6: model.JudgeVerdictAccepted,
	}

	for _, test := range tests {
		tempMiles := test.TimeUsage.Secs*1000 + test.TimeUsage.Nanos/1000000
		Priority = min(Priority, verdictPriorityMap[model.JudgeVerdict(test.Verdict)])
		AvgMiles += tempMiles
		MaxMiles = max(MaxMiles, tempMiles)
		AvgMemoryBytes += test.MemoryUsageBytes
		MaxMemoryBytes = max(MaxMemoryBytes, test.MemoryUsageBytes)
	}

	AvgMiles /= float64(totolTestPoint)
	AvgMemoryBytes /= float64(totolTestPoint)
	finalVerdict := priorityVerdictMap[Priority]

	core.AppLogger().Debugln(totolTestPoint, finalVerdict, AvgMiles, MaxMiles, AvgMemoryBytes, MaxMemoryBytes)

	// model.JudgeResult{
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
