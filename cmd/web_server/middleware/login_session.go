package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	auth_module "github.com/oj-lab/platform/modules/auth"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
)

const (
	loginSessionKeyIdCookieName      = "LS_KEY_ID"
	loginSessionKeyAccountCookieName = "LS_KEY_ACCOUNT"
	loginSessionGinCtxKey            = "login_session"
)

func BuildHandleRequireLoginWithRoles(roles []string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ls, err := GetLoginSessionFromGinCtx(ginCtx)
		if err != nil {
			gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
			ginCtx.Abort()
			return
		}
		ginCtx.Set(loginSessionGinCtxKey, ls)

		ginCtx.Next()
	}
}

func HandleRequireLogin(ginCtx *gin.Context) {
	ls, err := GetLoginSessionFromGinCtx(ginCtx)
	if err != nil {
		gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
		ginCtx.Abort()
		return
	}
	ginCtx.Set(loginSessionGinCtxKey, ls)

	ginCtx.Next()
}

func GetLoginSessionFromGinCtx(ginCtx *gin.Context) (*auth_module.LoginSession, error) {
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
	key := auth_module.LoginSessionKey{
		Account: lsAccount,
		Id:      lsId,
	}

	ls, err := auth_module.GetLoginSession(ginCtx, key)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func SetLoginSessionKeyCookie(ginCtx *gin.Context, key auth_module.LoginSessionKey) {
	ginCtx.SetCookie(loginSessionKeyAccountCookieName, key.Account,
		int(auth_module.LoginSessionDuration.Seconds()), "/", "", false, true)
	ginCtx.SetCookie(loginSessionKeyIdCookieName, key.Id.String(),
		int(auth_module.LoginSessionDuration.Seconds()), "/", "", false, true)
}
