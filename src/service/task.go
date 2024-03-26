package service

import (
	"context"
	"encoding/json"
	"fmt"

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

func parseVerdictJson(verdictString string) (model.JudgeVerdict, error) {
	var tests []map[string]interface{}
	err := json.Unmarshal([]byte(verdictString), &tests)
	if err != nil {
		return "", err
	}

	var (
		totolTestPoint = len(tests)
		Priority       = 6
		AvgTime        = 0.0
		MaxTime        = 0.0
		AvgMemory      = 0.0
		MaxMemory      = 0.0
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
		var (
			tmpVerdict = ""
			tmpNanos   = 0.0
			tmpSecs    = 0.0
			tmpMemory  = 0.0
		)
		tmpVerdict, ok := test["verdict"].(string)
		if !ok {
			return "", fmt.Errorf("verdict not exists")
		}
		tmpMemory, ok = test["memory_usage_bytes"].(float64)
		if !ok {
			return "", fmt.Errorf("memory_usage_bytes not exists")
		}
		tmpMap, ok := test["time_usage"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("time_usage not exists")
		}
		tmpNanos, ok = tmpMap["nanos"].(float64)
		if !ok {
			return "", fmt.Errorf("nanos not exists")
		}
		tmpSecs, ok = tmpMap["secs"].(float64)
		if !ok {
			return "", fmt.Errorf("secs not exists")
		}

		tmpMiles := tmpSecs*1000 + tmpNanos/1000000
		// core.AppLogger().Debugln(tmpVerdict, tmpNanos, tmpSecs, tmpMemory, tmpMiles)

		Priority = min(Priority, verdictPriorityMap[model.JudgeVerdict(tmpVerdict)])
		AvgTime += tmpMiles
		MaxTime = max(MaxTime, tmpMiles)
		AvgMemory += tmpMemory
		MaxMemory = max(MaxMemory, tmpMemory)
	}
	AvgTime /= float64(totolTestPoint)
	AvgMemory /= float64(totolTestPoint)
	finalVerdict := priorityVerdictMap[Priority]

	core.AppLogger().Debugln(totolTestPoint, finalVerdict, AvgTime, MaxTime, AvgMemory, MaxMemory)

	// model.JudgeResult{
	// 	MainVerdict:    finalVerdict,
	// 	TestPointCount: uint64(totolTestPoint),
	// 	MaxTimeMs:      uint64(MaxTime),
	// 	AverageTimeMs:  uint64(AvgTime),
	// 	maxTimeMs:      uint64(MaxTime),
	// 	AverageMemory:  uint64(AvgMemory),
	// 	MaxMemory:      uint64(MaxMemory),
	// }
	return finalVerdict, nil
}
