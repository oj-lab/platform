package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	config_module "github.com/oj-lab/platform/modules/config"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
)

const (
	internalTokenProp = "service.internal_token"
)

var (
	internalToken string
)

func init() {
	internalToken = config_module.AppConfig().GetString(internalTokenProp)
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
