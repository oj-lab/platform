package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/modules"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

func GetServiceError(ginErr gin.Error) *modules.SeviceError {
	if modules.IsServiceError(ginErr.Meta) {
		return ginErr.Meta.(*modules.SeviceError)
	} else {
		serviceErr := modules.NewInternalError(fmt.Sprintf("%v", ginErr.Err))
		serviceErr.CaptureStackTrace()
		return serviceErr
	}
}

func HandleError(ginCtx *gin.Context) {
	ginCtx.Next()

	errCount := len(ginCtx.Errors)
	if errCount > 0 {
		log.AppLogger().Errorf("Last error from GIN middleware: %+v", ginCtx.Errors[errCount-1].Err)
		err := GetServiceError(*ginCtx.Errors[errCount-1])
		ginCtx.JSON(err.Code, gin.H{
			"code": err.Code,
			"msg":  err.Msg,
		})
	}
}
