package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	"github.com/oj-lab/oj-lab-platform/modules"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
)

func SetupJudgeResultRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.PUT("/task/report/result-count", putReportJudgeResultCount)
		g.POST("/task/report/result", postReportJudgeResult)
	}
}

type ReportJudgeResultCountBody struct {
	JudgeUID    string `json:"judgeUID"`
	ResultCount uint   `json:"resultCount"`
}

func putReportJudgeResultCount(ginCtx *gin.Context) {
	body := ReportJudgeResultCountBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		modules.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
		return
	}

	judgeUID, err := uuid.Parse(body.JudgeUID)
	if err != nil {
		modules.NewInvalidParamError("judgeUID", "invalid judgeUID").AppendToGin(ginCtx)
		return
	}

	if err := judge_service.ReportJudgeResultCount(
		ginCtx, judgeUID, body.ResultCount,
	); err != nil {
		modules.NewInternalError(err.Error()).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}

type ReportJudgeResultBody struct {
	JudgeUIDString  string `json:"judgeUID"`
	VerdictString   string `json:"verdict"`
	TimeUsageMS     uint   `json:"timeUsageMS"`
	MemoryUsageByte uint   `json:"memoryUsageByte"`
}

func postReportJudgeResult(ginCtx *gin.Context) {
	body := ReportJudgeResultBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		modules.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
	}

	judgeUID, err := uuid.Parse(body.JudgeUIDString)
	if err != nil {
		modules.NewInvalidParamError("judgeUID", "invalid judgeUID").AppendToGin(ginCtx)
		return
	}
	verdict := judge_model.JudgeVerdict(body.VerdictString)
	if !verdict.IsValid() {
		modules.NewInvalidParamError("verdict", "invalid verdict").AppendToGin(ginCtx)
		return
	}

	judge_service.CreateJudgeResult(ginCtx, judge_model.JudgeResult{
		JudgeUID:        judgeUID,
		Verdict:         verdict,
		TimeUsageMS:     body.TimeUsageMS,
		MemoryUsageByte: body.MemoryUsageByte,
	})
}
