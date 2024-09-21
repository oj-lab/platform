package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oj-lab/platform/cmd/web_server/middleware"
	"github.com/oj-lab/platform/models"
	judge_model "github.com/oj-lab/platform/models/judge"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
	judge_service "github.com/oj-lab/platform/services/judge"
	user_service "github.com/oj-lab/platform/services/user"
)

func SetupJudgeRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.GET("", middleware.HandleRequireLogin, getJudgeList)
		g.GET("/:uid", middleware.HandleRequireLogin, getJudge)
	}
}

func getJudge(ginCtx *gin.Context) {
	uidString := ginCtx.Param("uid")
	uid, err := uuid.Parse(uidString)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "uid", "invalid uid")
		return
	}

	judge, err := judge_service.GetJudge(ginCtx, uid)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get judge: %v", err))
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
	limit, err := gin_utils.QueryInt(ginCtx, "limit", 10)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "limit", err.Error())
		return
	}
	offset, err := gin_utils.QueryInt(ginCtx, "offset", 0)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "offset", err.Error())
		return
	}
	selfOnly := gin_utils.QueryBool(ginCtx, "self_only", false)

	options := judge_model.GetJudgeOptions{
		Limit:          &limit,
		Offset:         &offset,
		OrderByColumns: []models.OrderByColumnOption{{Column: "create_at", Desc: true}},
	}
	if selfOnly {
		ls, err := middleware.GetLoginSessionFromGinCtx(ginCtx)
		if err != nil {
			gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
			return
		}
		user, err := user_service.GetUser(ginCtx, ls.Key.Account)
		if err != nil {
			gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get user: %v", err))
			return
		}
		options.UserAccount = user.Account
	}

	judges, total, err := judge_service.GetJudgeList(ginCtx, options)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get judge list: %v", err))
		return
	}

	ginCtx.JSON(200, getJudgeListResponse{
		Total: total,
		List:  judges,
	})
}
