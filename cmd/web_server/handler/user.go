package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/cmd/web_server/middleware"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	"github.com/oj-lab/oj-lab-platform/modules"
	user_service "github.com/oj-lab/oj-lab-platform/services/user"
)

func SetupUserRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/user")
	{
		g.PUT("", updateUser)
		g.GET("/health/*any", func(ginCtx *gin.Context) {
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

	ls, err := user_service.StartLoginSession(ginCtx, body.Account, body.Password)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to login: %v", err)).AppendToGin(ginCtx)
		return
	}
	middleware.SetLoginSessionKeyCookie(ginCtx, ls.Key)

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
	ls, err := middleware.GetLoginSessionFromGinCtx(ginCtx)
	if err != nil {
		modules.NewUnauthorizedError("cannot load login session from cookie").AppendToGin(ginCtx)
		return
	}
	user, err := user_service.GetUser(ginCtx, ls.Key.Account)
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

type updateUserBody struct {
	User user_model.User `json:"user"`
}

func updateUser(ginCtx *gin.Context) {
	body := &updateUserBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		modules.NewInvalidParamError("body", err.Error()).AppendToGin(ginCtx)
		return
	}

	err = user_service.UpdateUser(ginCtx, body.User)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to update user: %v", err)).AppendToGin(ginCtx)
		return
	}
}
