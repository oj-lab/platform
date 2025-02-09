package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	core_module "github.com/oj-lab/platform/modules/core"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
)

const (
	internalTokenConfigKey = "service.internal_token"
)

var (
	internalToken string
)

func init() {
	internalToken = core_module.Config.GetString(internalTokenConfigKey)
}

func HandleRequireInternalToken(ginCtx *gin.Context) {
	incommingToken := ginCtx.GetHeader("Authorization")
	if fmt.Sprintf("Bearer %s", internalToken) != incommingToken {
		gin_utils.NewUnauthorizedError(ginCtx, "invalid internal token")
		ginCtx.Abort()
		return
	}

	ginCtx.Next()
}
