package handler

import (
	"github.com/gin-gonic/gin"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
	"github.com/redis/go-redis/v9"
)

func SetupJudgeRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.POST("/task/pick", postPickJudgeTask)
		g.POST("/task/report", postReportJudgeTaskResult)
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

type ReportJudgeTaskResultBody struct {
	Consumer    string `json:"consumer"`
	StreamID    string `json:"stream_id"`
	VerdictJson string `json:"verdict_json"`
}

func postReportJudgeTaskResult(ginCtx *gin.Context) {
	body := ReportJudgeTaskResultBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	if err := judge_service.ReportJudgeTaskResult(ginCtx, body.Consumer, body.StreamID, body.VerdictJson); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}
