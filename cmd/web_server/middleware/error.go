package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log_module "github.com/oj-lab/platform/modules/log"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
)

func GetServiceError(ginErr gin.Error) *gin_utils.SeviceError {
	if gin_utils.IsServiceError(ginErr.Meta) {
		return ginErr.Meta.(*gin_utils.SeviceError)
	} else {
		serviceErr := gin_utils.SeviceError{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("%v", ginErr.Err),
		}
		serviceErr.CaptureStackTrace()
		return &serviceErr
	}
}

func HandleError(ginCtx *gin.Context) {
	ginCtx.Next()

	errCount := len(ginCtx.Errors)
	if errCount > 0 {
		log_module.AppLogger().Errorf("Last error from GIN middleware: %+v", ginCtx.Errors[errCount-1].Err)
		err := GetServiceError(*ginCtx.Errors[errCount-1])
		ginCtx.JSON(err.Code, gin.H{
			"code": err.Code,
			"msg":  err.Msg,
		})
	}
}
