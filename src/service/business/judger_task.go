// ! Deprected

package business

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OJ-lab/oj-lab-services/src/core"
	asynqAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/asynq"
	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	judgerAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/judger"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MuxPatternJudger               = "judger"
	TaskNameJudgerTrackAllState    = "judger:track_all_state"
	TaskNameJudgerGetState         = "judger:get_state"
	TaskNameJudgerHandleSubmission = "judger:handle_submission"
)

func NewTaskJudgerTrackAllState() *asynq.Task {
	return asynq.NewTask(TaskNameJudgerTrackAllState, nil)
}

func NewTaskJudgerGetState(judger model.Judger) *asynq.Task {
	judgerJson, err := json.Marshal(judger)
	if err != nil {
		panic(err)
	}
	return asynq.NewTask(TaskNameJudgerGetState, judgerJson)
}

func GetAsynqMuxJudger() asynqAgent.AsynqMux {
	serveMux := asynq.NewServeMux()
	serveMux.Use(asynqAgent.LoggingMiddleware)
	serveMux.HandleFunc(TaskNameJudgerTrackAllState, handleTaskJudgerTrackAllState)
	serveMux.HandleFunc(TaskNameJudgerGetState, handleTaskJudgerGetState)
	serveMux.HandleFunc(TaskNameJudgerHandleSubmission, handleTaskJudgerHandleSubmission)

	return asynqAgent.AsynqMux{
		Pattern:  MuxPatternJudger,
		ServeMux: serveMux,
	}
}

func handleTaskJudgerTrackAllState(ctx context.Context, task *asynq.Task) error {
	db := gormAgent.GetDefaultDB()
	judgerList, err := mapper.GetJudgerList(db)
	if err != nil {
		return err
	}
	core.AppLogger().Infof("judger list: %v", judgerList)

	asynqClient := asynqAgent.GetDefaultTaskClient()
	for _, judger := range judgerList {
		err := asynqClient.EnqueueTask(
			NewTaskJudgerGetState(judger),
			asynq.TaskID(fmt.Sprintf("%s:%s", TaskNameJudgerGetState, judger.Host)),
		)
		if err != nil {
			core.AppLogger().Errorf("failed to enqueue task %s: %v", TaskNameJudgerGetState, err)
		}
	}

	return nil
}

func handleTaskJudgerGetState(ctx context.Context, task *asynq.Task) error {
	db := gormAgent.GetDefaultDB()
	var judger model.Judger
	if err := json.Unmarshal(task.Payload(), &judger); err != nil {
		return err
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("host = ?", judger.Host).First(&judger).Error; err != nil {
			return err
		}

		judgerClient := judgerAgent.JudgerClient{
			Host: judger.Host,
		}
		judgerStateString, err := judgerClient.GetState()
		if err != nil {
			return err
		}
		judgerState := model.StringToJudgerState(judgerStateString)
		core.AppLogger().Debugf("Get Judger %v state=%v", judgerClient.Host, judgerState)

		if judgerState == model.JudgerStateIdle {
			core.AppLogger().Debugf("Judger %v is idle, try find submission to handle", judgerClient.Host)
			asynqClient := asynqAgent.GetDefaultTaskClient()
			err := asynqClient.EnqueueTask(
				NewTaskJudgerHandleSubmission(judger),
				asynq.TaskID(fmt.Sprintf("%s:%s", TaskNameJudgerHandleSubmission, judger.Host)),
			)
			if err != nil {
				core.AppLogger().Errorf("failed to enqueue task %s: %v", TaskNameJudgerHandleSubmission, err)
			}
		}

		if !judger.State.CanUpdate(model.JudgerStateIdle) {
			core.AppLogger().Debugf("Judger state is invalid or don't need update, ignoring this state update")
			return nil
		}
		judger.State = judgerState

		err = tx.Model(&judger).Update("state", judgerState).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func NewTaskJudgerHandleSubmission(judger model.Judger) *asynq.Task {
	judgerJson, err := json.Marshal(judger)
	if err != nil {
		panic(err)
	}
	return asynq.NewTask(TaskNameJudgerHandleSubmission, judgerJson)
}

func handleTaskJudgerHandleSubmission(ctx context.Context, task *asynq.Task) error {
	db := gormAgent.GetDefaultDB()
	var judger model.Judger
	if err := json.Unmarshal(task.Payload(), &judger); err != nil {
		return err
	}

	judgerClient := judgerAgent.JudgerClient{
		Host: judger.Host,
	}
	judgerStateString, err := judgerClient.GetState()
	if err != nil {
		return err
	}
	judgerState := model.StringToJudgerState(judgerStateString)
	if judgerState != model.JudgerStateIdle {
		core.AppLogger().Debugf("Judger %v is not idle, ignoring this task", judgerClient.Host)
		return nil
	}

	var submission model.JudgeTaskSubmission
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := mapper.BuildGetSubmissionTXByOptions(tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}), mapper.GetSubmissionOptions{
			Statuses: []model.SubmissionStatus{
				model.SubmissionStatusPending,
			},
			Limit: func() *int {
				limit := 1
				return &limit
			}(),
			OrderByColumns: []model.OrderByColumnOption{
				{
					Column: "create_at",
					Desc:   true,
				},
			},
		}, false).First(&submission).Error; err != nil {
			core.AppLogger().Errorf("Failed to get submission: %v", err)
			return err
		}
		core.AppLogger().Debugf("Get submission %v", submission.UID)
		submission.Status = model.SubmissionStatusRunning
		err = tx.Model(&submission).Update("status", submission.Status).Error
		if err != nil {
			core.AppLogger().Errorf("Failed to update submission status: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		core.AppLogger().Errorf("Failed to select submission: %v", err)
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("host = ?", judger.Host).First(&judger).Error; err != nil {
			return err
		}
		judger.State = model.JudgerStateBusy
		err = tx.Model(&judger).Update("state", judger.State).Error
		if err != nil {
			return err
		}

		core.AppLogger().Debugf("Judger %v start to handle submission %v", judgerClient.Host, submission.UID)
		judgeVerdict, err := judgerClient.PostJudgeSync(
			submission.ProblemSlug,
			judgerAgent.JudgeRequest{
				Language: string(submission.Language),
				Code:     submission.Code,
			},
		)
		if err != nil {
			core.AppLogger().Errorf("Failed to judge submission: %v", err)
			return err
		}
		core.AppLogger().Debugf("Get Judger %v verdict=%v", judgerClient.Host, judgeVerdict)

		submission.Status = model.SubmissionStatusFinished
		verdictBytes, err := json.Marshal(judgeVerdict)
		if err != nil {
			core.AppLogger().Errorf("Failed to marshal verdict: %v", err)
			return err
		}
		err = tx.Model(&submission).Updates(map[string]interface{}{
			"status":       submission.Status,
			"verdict_json": string(verdictBytes),
		}).Error
		if err != nil {
			return err
		}

		judger.State = model.JudgerStateIdle
		err = tx.Model(&judger).Update("state", judger.State).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		core.AppLogger().Errorf("Failed to handle submission: %v", err)
		err = db.Model(&submission).Update("status", model.SubmissionStatusPending).Error
		if err != nil {
			core.AppLogger().Errorf("failed to update submission status: %v", err)
		}

		return err
	}
	return nil
}
