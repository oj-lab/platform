package handler

import (
	"github.com/gin-gonic/gin"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
	judge_service "github.com/oj-lab/platform/services/judge"
)

func SetupRankRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/rank")
	{
		g.GET("", getRankList)
		// g.GET("/:account", getUserRank)
	}
}

// getUserRank
//
//	@Router			/rank/{account} [get]
//	@Summary		Get a user rank
//	@Description	Get a rank for user
//	@Tags			problem
//	@Accept			json
//	@Success		200
// func getUserRank(ginCtx *gin.Context) {

// }

// getRankList
//
//	@Router			/rank [get]
//	@Summary		Get rank list
//	@Description	Get rank list
//	@Tags			rank
//	@Accept			json
//	@Success		200
func getRankList(ginCtx *gin.Context) {
	limit, err := gin_utils.QueryInt(ginCtx, "limit", 100)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "limit", err.Error())
		return
	}
	offset, err := gin_utils.QueryInt(ginCtx, "offset", 0)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "offset", err.Error())
		return
	}

	rankInfoList, total, err := judge_service.GetRankList(
		ginCtx,
		nil,
		&limit, &offset,
	)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, err.Error())
		return
	}

	ginCtx.JSON(200, gin.H{
		"total": total,
		"list":  rankInfoList,
	})
}
