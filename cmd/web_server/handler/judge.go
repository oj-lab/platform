package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/models"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	"github.com/oj-lab/oj-lab-platform/modules"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
)

func SetupJudgeRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.GET("", getJudgeList)
		g.GET("/:uid", getJudge)
	}
}

func getJudge(ginCtx *gin.Context) {
	uidString := ginCtx.Param("uid")
	uid, err := uuid.Parse(uidString)
	if err != nil {
		modules.NewInvalidParamError("uid", "invalid uid").AppendToGin(ginCtx)
		return
	}

	judge, err := judge_service.GetJudge(ginCtx, uid)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get judge: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, judge)
}

type getJudgeListResponse struct {
	Total int64                `json:"total"`
	List  []*judge_model.Judge `json:"list"`
}

// Get Judge List
//
//	@Summary		Get Judge list
//	@Description	Get Judge list
//	@Tags			judge
//	@Accept			json
//	@Param			limit	query	int	false	"limit"
//	@Param			offset	query	int	false	"offset"
//	@Router			/judge [get] getJudgeListResponse
func getJudgeList(ginCtx *gin.Context) {
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

	options := judge_model.GetJudgeOptions{
		Limit:          &limit,
		Offset:         &offset,
		OrderByColumns: []models.OrderByColumnOption{{Column: "create_at", Desc: true}},
	}

	judges, total, err := judge_service.GetJudgeList(ginCtx, options)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get judge list: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, getJudgeListResponse{
		Total: total,
		List:  judges,
	})
}
