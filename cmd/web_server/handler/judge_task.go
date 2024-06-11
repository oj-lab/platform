package handler

import (
	"github.com/gin-gonic/gin"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	"github.com/oj-lab/oj-lab-platform/modules"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
	"github.com/redis/go-redis/v9"
)

func SetupJudgeTaskRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.POST("/task/pick", postPickJudgeTask)
		g.PUT("/task/report", putReportJudgeTask)
	}
}

type PickJudgeTaskBody struct {
	Consumer string `json:"consumer"`
}

func postPickJudgeTask(ginCtx *gin.Context) {
	body := PickJudgeTaskBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	task, err := judge_service.PickJudgeTask(ginCtx, body.Consumer)
	if err == redis.Nil {
		ginCtx.Status(204)
		return
	}

	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"task": task,
	})
}

type ReportJudgeTaskBody struct {
	Consumer      string `json:"consumer"`
	RedisStreamID string `json:"redisStreamID"`
	VerdictString string `json:"verdict"`
}

func putReportJudgeTask(ginCtx *gin.Context) {
	body := ReportJudgeTaskBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	verdict := judge_model.JudgeVerdict(body.VerdictString)
	if !verdict.IsValid() {
		modules.NewInvalidParamError("verdict", "invalid verdict").AppendToGin(ginCtx)
		return
	}

	if err := judge_service.ReportJudgeTask(
		ginCtx, body.Consumer, body.RedisStreamID, verdict,
	); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}
