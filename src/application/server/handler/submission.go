package handler

import (
	"strconv"

	"github.com/OJ-lab/oj-lab-services/src/core"
	"github.com/OJ-lab/oj-lab-services/src/service"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"github.com/gin-gonic/gin"
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

	submission, svcErr := service.GetJudgeTaskSubmission(ginCtx, uid)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, gin.H{
		"submission": submission,
	})
}

type getSubmissionListResponse struct {
	Total int64                        `json:"total"`
	List  []*model.JudgeTaskSubmission `json:"list"`
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
		core.NewInvalidParamError("limit", "invalid limit").AppendToGin(ginCtx)
		return
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		core.NewInvalidParamError("offset", "invalid offset").AppendToGin(ginCtx)
		return
	}

	options := mapper.GetSubmissionOptions{
		Limit:  &limit,
		Offset: &offset,
	}

	submissions, total, svcErr := service.GetJudgeTaskSubmissionList(ginCtx, options)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, getSubmissionListResponse{
		Total: total,
		List:  submissions,
	})
}
