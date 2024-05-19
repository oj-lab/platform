package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/modules"
	"github.com/oj-lab/oj-lab-platform/modules/auth"
)

const (
	loginSessionCookieMaxAge = time.Hour * 24 * 7
	loginSessionIdCookieName = "LS_ID"
	loginSessionGinCtxKey    = "login_session"
)

func HandleRequireLogin(ginCtx *gin.Context) {
	cookie, err := ginCtx.Cookie(loginSessionIdCookieName)
	if err != nil {
		modules.NewUnauthorizedError("login session not found").AppendToGin(ginCtx)
		ginCtx.Abort()
		return
	}

	ls, err := auth.CheckLoginSession(ginCtx, cookie)
	if err != nil {
		modules.NewUnauthorizedError("invalid login session").AppendToGin(ginCtx)
		ginCtx.Abort()
		return
	}

	ginCtx.Set(loginSessionGinCtxKey, ls)

	ginCtx.Next()
}

func GetLoginSession(ginCtx *gin.Context) *auth.LoginSession {
	ls, exist := ginCtx.Get(loginSessionGinCtxKey)
	if !exist {
		return nil
	}
	return ls.(*auth.LoginSession)
}

func SetLoginSessionCookie(ginCtx *gin.Context, lsId string) {
	ginCtx.SetCookie(loginSessionIdCookieName, lsId,
		int(loginSessionCookieMaxAge.Seconds()), "/", "", false, true)
}
