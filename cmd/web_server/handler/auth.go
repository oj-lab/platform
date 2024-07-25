package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/cmd/web_server/middleware"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	"github.com/oj-lab/oj-lab-platform/modules"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	auth_module "github.com/oj-lab/oj-lab-platform/modules/auth"
	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
	user_service "github.com/oj-lab/oj-lab-platform/services/user"
)

func SetupOauthRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/auth")
	{
		g.GET("/github/callback", githubCallback)
		g.Any("/github", loginGithub)

		g.POST("/password", loginByPassword)
	}
}

func githubCallback(ginCtx *gin.Context) {
	code := ginCtx.Query("code")
	tokenResponse, err := auth_module.GetGithubAccessToken(code)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	log_module.AppLogger().WithField("tokenResponse", tokenResponse).Info("github callback")

	ginCtx.JSON(200, nil)
}

func loginGithub(ginCtx *gin.Context) {
	u, err := auth_module.GetGithubOauthEntryURL()
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get github oauth entry url: %v", err)).AppendToGin(ginCtx)
	}
	ginCtx.Redirect(http.StatusFound, u.String())
}

type loginBody struct {
	Account  string `json:"account" example:"admin"`
	Password string `json:"password" example:"admin"`
}

// LoginByPassword
//
//	@Summary		Login by account and password
//	@Description	A Cookie will be set if login successfully
//	@Tags			user
//	@Accept			json
//	@Param			loginBody	body	loginBody	true	"body"
//	@Router			/user/login [post]
//	@Success		200
func loginByPassword(ginCtx *gin.Context) {
	body := &loginBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		modules.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
		return
	}

	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUserByAccountPassword(db, body.Account, body.Password)
	if err != nil {
		modules.NewUnauthorizedError("account or password incorrect").AppendToGin(ginCtx)
	}

	ls, err := user_service.StartLoginSession(ginCtx, user.Account)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to login: %v", err)).AppendToGin(ginCtx)
		return
	}
	middleware.SetLoginSessionKeyCookie(ginCtx, ls.Key)

	ginCtx.Status(http.StatusOK)
}
