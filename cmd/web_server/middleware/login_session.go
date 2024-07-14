package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/modules"
	"github.com/oj-lab/oj-lab-platform/modules/auth"
)

const (
	loginSessionCookieMaxAge         = time.Hour * 24 * 7
	loginSessionKeyIdCookieName      = "LS_KEY_ID"
	loginSessionKeyAccountCookieName = "LS_KEY_ACCOUNT"
	loginSessionGinCtxKey            = "login_session"
)

func BuildHandleRequireLoginWithRoles(roles []string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ls, err := GetLoginSessionFromGinCtx(ginCtx)
		if err != nil {
			modules.NewUnauthorizedError("cannot load login session from cookie").AppendToGin(ginCtx)
			ginCtx.Abort()
			return
		}
		ginCtx.Set(loginSessionGinCtxKey, ls)

		if !ls.Data.HasRoles(roles) {
			modules.NewUnauthorizedError(fmt.Sprintf("require roles: %v", roles)).AppendToGin(ginCtx)
			ginCtx.Abort()
			return
		}

		ginCtx.Next()
	}
}

func HandleRequireLogin(ginCtx *gin.Context) {
	ls, err := GetLoginSessionFromGinCtx(ginCtx)
	if err != nil {
		modules.NewUnauthorizedError("cannot load login session from cookie").AppendToGin(ginCtx)
		ginCtx.Abort()
		return
	}
	ginCtx.Set(loginSessionGinCtxKey, ls)

	ginCtx.Next()
}

func GetLoginSessionFromGinCtx(ginCtx *gin.Context) (*auth.LoginSession, error) {
	lsAccount, err := ginCtx.Cookie(loginSessionKeyAccountCookieName)
	if err != nil {
		return nil, err
	}
	lsIdString, err := ginCtx.Cookie(loginSessionKeyIdCookieName)
	if err != nil {
		return nil, err
	}
	lsId, err := uuid.Parse(lsIdString)
	if err != nil {
		return nil, err
	}
	key := auth.LoginSessionKey{
		Account: lsAccount,
		Id:      lsId,
	}

	ls, err := auth.GetLoginSession(ginCtx, key)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func SetLoginSessionKeyCookie(ginCtx *gin.Context, key auth.LoginSessionKey) {
	ginCtx.SetCookie(loginSessionKeyAccountCookieName, key.Account,
		int(loginSessionCookieMaxAge.Seconds()), "/", "", false, true)
	ginCtx.SetCookie(loginSessionKeyIdCookieName, key.Id.String(),
		int(loginSessionCookieMaxAge.Seconds()), "/", "", false, true)
}
