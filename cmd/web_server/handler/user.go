package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/modules"
	"github.com/oj-lab/oj-lab-platform/modules/middleware"
	user_service "github.com/oj-lab/oj-lab-platform/services/user"
)

func SetupUserRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/user")
	{
		g.GET("/health", func(ginCtx *gin.Context) {
			ginCtx.String(http.StatusOK, "Hello, this is user service")
		})
		g.POST("/login", login)
		g.GET("/me", middleware.HandleRequireLogin, me)
		g.GET("/check-exist", checkUserExist)
	}
}

type loginBody struct {
	Account  string `json:"account" example:"admin"`
	Password string `json:"password" example:"admin"`
}

// Login
//
//	@Summary		Login by account and password
//	@Description	A Cookie will be set if login successfully
//	@Tags			user
//	@Accept			json
//	@Param			loginBody	body	loginBody	true	"body"
//	@Router			/user/login [post]
//	@Success		200
func login(ginCtx *gin.Context) {
	body := &loginBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		modules.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
		return
	}

	lsId, err := user_service.StartLoginSession(ginCtx, body.Account, body.Password)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to login: %v", err)).AppendToGin(ginCtx)
		return
	}
	middleware.SetLoginSessionCookie(ginCtx, lsId.String())

	ginCtx.Status(http.StatusOK)
}

// Me
//
//	@Summary		Get current user
//	@Description	If correctly logined with cookie, return current user
//	@Tags			user
//	@Router			/user/me [get]
//	@Success		200
//	@Failure		401
func me(ginCtx *gin.Context) {
	ls := middleware.GetLoginSession(ginCtx)
	if ls == nil {
		modules.NewUnauthorizedError("not logined").AppendToGin(ginCtx)
		return
	}
	user, err := user_service.GetUser(ginCtx, ls.Account)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get user: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(http.StatusOK, user)
}

func checkUserExist(ginCtx *gin.Context) {
	account := ginCtx.Query("account")
	if account == "" {
		modules.NewInvalidParamError("account", "account cannot be empty").AppendToGin(ginCtx)
		return
	}

	exist, err := user_service.CheckUserExist(ginCtx, account)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to check user exist: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{
		"exist": exist,
	})
}
