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

func QueryString(ginCtx *gin.Context, key string, defaultValue string) string {
	ValueStr := ginCtx.Query(key)
	if ValueStr == "" {
		return defaultValue
	}
	return ValueStr
}

func QueryBool(ginCtx *gin.Context, key string, defaultValue bool) bool {
	ValueStr := ginCtx.Query(key)
	switch ValueStr {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}
