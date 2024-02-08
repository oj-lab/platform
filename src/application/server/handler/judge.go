package handler

import (
	"github.com/OJ-lab/oj-lab-services/src/service"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupJudgeRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.POST("/add-judger", postJudger)
		g.POST("/task/pick", postPickJudgeTask)
		g.POST("/task/report", postReportJudgeTaskResult)
	}
}

func postJudger(ginCtx *gin.Context) {
	judger := model.Judger{}
	if err := ginCtx.ShouldBindJSON(&judger); err != nil {
		ginCtx.Error(err)
		return
	}

	if err := service.AddJudger(ginCtx, judger); err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}

type PickJudgeTaskBody struct {
	Consumer string `json:"consumer"`
}

func postPickJudgeTask(ginCtx *gin.Context) {
	body := PickJudgeTaskBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		ginCtx.Error(err)
		return
	}

	task, err := service.PickJudgeTask(ginCtx, body.Consumer)
	if err == redis.Nil {
		ginCtx.Status(204)
		return
	}

	if err != nil {
		ginCtx.Error(err)
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
		ginCtx.Error(err)
		return
	}

	if err := service.ReportJudgeTaskResult(ginCtx, body.Consumer, body.StreamID, body.VerdictJson); err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}
