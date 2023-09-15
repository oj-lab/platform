package middleware

import (
	"github.com/OJ-lab/oj-lab-services/package/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleError(ginCtx *gin.Context) {
	ginCtx.Next()

	errCount := len(ginCtx.Errors)
	if errCount > 0 {
		logrus.Errorf("Last error from GIN middleware: %+v", ginCtx.Errors[errCount-1].Err)
		err := core.GetServiceErrorFromGinError(*ginCtx.Errors[errCount-1])
		ginCtx.JSON(err.Code, gin.H{
			"code": err.Code,
			"msg":  err.Msg,
		})
	}
}
