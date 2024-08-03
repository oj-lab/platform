package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gin_utils "github.com/oj-lab/oj-lab-platform/modules/utils/gin"
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
		gin_utils.NewInvalidParamError(ginCtx, "body", "invalid body")
		return
	}

	judgeUID, err := uuid.Parse(body.JudgeUID)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "judgeUID", "invalid judgeUID")
		return
	}

	if err := judge_service.ReportJudgeResultCount(
		ginCtx, judgeUID, body.ResultCount,
	); err != nil {
		gin_utils.NewInternalError(ginCtx, err.Error())
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
		gin_utils.NewInvalidParamError(ginCtx, "body", "invalid body")
	}

	judgeUID, err := uuid.Parse(body.JudgeUIDString)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "judgeUID", "invalid judgeUID")
		return
	}
	verdict := judge_model.JudgeVerdict(body.VerdictString)
	if !verdict.IsValid() {
		gin_utils.NewInvalidParamError(ginCtx, "verdict", "invalid verdict")
		return
	}

	_, err = judge_service.CreateJudgeResult(ginCtx, judge_model.JudgeResult{
		JudgeUID:        judgeUID,
		Verdict:         verdict,
		TimeUsageMS:     body.TimeUsageMS,
		MemoryUsageByte: body.MemoryUsageByte,
	})
	if err != nil {
		gin_utils.NewInternalError(ginCtx, err.Error())
		return
	}
}
