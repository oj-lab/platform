package modules

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type SeviceError struct {
	Code       int       `json:"code"`
	Msg        string    `json:"msg"`
	stackTrace []uintptr `json:"-"`
}

func (se *SeviceError) ToGinError() *gin.Error {
	return &gin.Error{
		Err:  fmt.Errorf("%v", se.Msg),
		Type: gin.ErrorTypePrivate,
		Meta: se,
	}
}

func (se *SeviceError) AppendToGin(ginCtx *gin.Context) {
	ginCtx.Errors = append(ginCtx.Errors, se.ToGinError())
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

func NewInternalError(msg string) *SeviceError {
	return &SeviceError{
		Code: http.StatusInternalServerError,
		Msg:  msg,
	}
}

func NewUnauthorizedError(msg string) *SeviceError {
	return &SeviceError{
		Code: http.StatusUnauthorized,
		Msg:  msg,
	}
}

func NewInvalidParamError(param string, hints ...string) *SeviceError {
	msg := fmt.Sprintf("invalid param: %s", param)
	for _, hint := range hints {
		msg += fmt.Sprintf(", %s", hint)
	}

	return &SeviceError{
		Code: http.StatusBadRequest,
		Msg:  msg,
	}
}
