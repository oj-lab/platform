package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/models"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	"github.com/oj-lab/oj-lab-platform/modules"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
)

func SetupSubmissionRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/submission")
	{
		g.GET("", getSubmissionList)
		g.GET("/:uid", getSubmission)
	}
}

func getSubmission(ginCtx *gin.Context) {
	uid := ginCtx.Param("uid")

	submission, err := judge_service.GetJudgeTaskSubmission(ginCtx, uid)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get submission: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, gin.H{
		"UID":           submission.UID,
		"redisStreamID": submission.RedisStreamID,
		"userAccount":   submission.UserAccount,
		"user":          submission.User, // include User metadata, If is needed
		"problemSlug":   submission.ProblemSlug,
		"problem":       submission.Problem, // include Problem metadata, If is needed
		"code":          submission.Code,
		"language":      submission.Language,
		"status":        submission.Status,
		"verdictJson":   submission.VerdictJson,
		"mainResult":    submission.MainResult,
	})
}

type getSubmissionListResponse struct {
	Total int64                              `json:"total"`
	List  []*judge_model.JudgeTaskSubmission `json:"list"`
}

// Get Submission List
//
//	@Summary		Get submission list
//	@Description	Get submission list
//	@Tags			submission
//	@Accept			json
//	@Param			limit	query	int	false	"limit"
//	@Param			offset	query	int	false	"offset"
//	@Router			/submission [get] getSubmissionListResponse
func getSubmissionList(ginCtx *gin.Context) {
	limitQuery, _ := ginCtx.GetQuery("limit")
	offsetQuery, _ := ginCtx.GetQuery("offset")
	if limitQuery == "" {
		limitQuery = "10"
	}
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		modules.NewInvalidParamError("limit", "invalid limit").AppendToGin(ginCtx)
		return
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		modules.NewInvalidParamError("offset", "invalid offset").AppendToGin(ginCtx)
		return
	}

	options := judge_model.GetSubmissionOptions{
		Limit:          &limit,
		Offset:         &offset,
		OrderByColumns: []models.OrderByColumnOption{{Column: "create_at", Desc: true}},
	}

	submissions, total, err := judge_service.GetJudgeTaskSubmissionList(ginCtx, options)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get submission list: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, getSubmissionListResponse{
		Total: total,
		List:  submissions,
	})
}
