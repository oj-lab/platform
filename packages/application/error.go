package application

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SeviceError struct {
	Code       int       `json:"code"`
	Msg        string    `json:"msg"`
	stackTrace []uintptr `json:"-"`
}

func (se *SeviceError) CaptureStackTrace() *SeviceError {
	se.stackTrace = []uintptr{}
	runtime.Callers(2, se.stackTrace)

	return se
}

func IsServiceError(err interface{}) bool {
	_, ok := err.(*SeviceError)
	return ok
}

func WrapToServiceError(err interface{}) *SeviceError {
	if IsServiceError(err) {
		return err.(*SeviceError)
	} else {
		serviceErr := NewInternalError(fmt.Sprintf("%v", err))
		serviceErr.CaptureStackTrace()
		return serviceErr
	}
}

func NewInternalError(msg string) *SeviceError {
	return &SeviceError{
		Code: 500,
		Msg:  msg,
	}
}

func NewUnAuthorizedError(msg string) *SeviceError {
	return &SeviceError{
		Code: 401,
		Msg:  msg,
	}
}

func HandleError(ginCtx *gin.Context) {
	ginCtx.Next()

	errCount := len(ginCtx.Errors)
	if errCount > 0 {
		logrus.Errorf("Last error from GIN middleware: %+v", ginCtx.Errors[errCount-1].Err)
		err := WrapToServiceError(ginCtx.Errors[errCount-1].Err)
		ginCtx.JSON(err.Code, gin.H{
			"code": err.Code,
			"msg":  err.Msg,
		})
	}
}
