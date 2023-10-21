package business

import (
	"context"
	"encoding/json"

	"github.com/OJ-lab/oj-lab-services/core"
	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	judgerAgent "github.com/OJ-lab/oj-lab-services/core/agent/judger"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MuxPatternJudger            = "judger"
	TaskNameJudgerTrackAllState = "judger:track_all_state"
	TaskNameJudgerGetState      = "judger:get_state"
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
	serveMux.HandleFunc(TaskNameJudgerTrackAllState, handleTaskJudgerTrackAllState)
	serveMux.HandleFunc(TaskNameJudgerGetState, handleTaskJudgerGetState)

	return asynqAgent.AsynqMux{
		Pattern:  MuxPatternJudger,
		ServeMux: serveMux,
	}
}

func handleTaskJudgerTrackAllState(ctx context.Context, task *asynq.Task) error {
	core.GetAppLogger().Info("handleTaskJudgerTrackAllState")
	db := gormAgent.GetDefaultDB()
	judgerList, err := mapper.GetJudgerList(db)
	if err != nil {
		return err
	}
	core.GetAppLogger().Infof("judger list: %v", judgerList)

	asynqClient := asynqAgent.GetDefaultTaskClient()
	for _, judger := range judgerList {
		asynqClient.EnqueueTask(
			NewTaskJudgerGetState(judger),
			asynq.TaskID(judger.Host),
		)
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
		core.GetAppLogger().Debugf("Get Judger %v state=%v", judgerClient.Host, judgerState)

		if !judger.State.CanUpdate(model.JudgerStateIdle) {
			core.GetAppLogger().Debugf("Judger state is invalid, ignoring this state update")
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
	core.GetAppLogger().Debugf("Successfully handled task %s", task.Type())
	return nil
}
