package handler

import (
	"strconv"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/gin-gonic/gin"
)

func SetupSubmissionRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/submission")
	{
		g.GET("", getSubmissionList)
	}
}

type getSubmissionResponseBody struct {
	Total int64                        `json:"total"`
	Data  []*model.JudgeTaskSubmission `json:"data"`
}

// Get Submission List
//
//	@Summary		Get submission list
//	@Description	Get submission list
//	@Tags			submission
//	@Accept			json
//	@Param			limit	query	int	false	"limit"
//	@Param			offset	query	int	false	"offset"
//	@Router			/submission [get]
//	@Success		200	{object}	getSubmissionResponseBody
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

	ginCtx.JSON(200, getSubmissionResponseBody{
		Total: total,
		Data:  submissions,
	})
}
