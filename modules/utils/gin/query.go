package gin_utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInt(ginCtx *gin.Context, key string, defaultValue int) (int, error) {
	ValueStr := ginCtx.Query(key)
	if ValueStr == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(ValueStr)
}
