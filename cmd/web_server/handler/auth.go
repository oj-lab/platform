package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/platform/cmd/web_server/middleware"
	user_model "github.com/oj-lab/platform/models/user"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	auth_module "github.com/oj-lab/platform/modules/auth"
	log_module "github.com/oj-lab/platform/modules/log"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
	user_service "github.com/oj-lab/platform/services/user"
)

const callbackURL = "/auth/github/callback"

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
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get github access token: %v", err))
		return
	}

	log_module.AppLogger().WithField("tokenResponse", tokenResponse).Info("github callback")
	githubUser, err := auth_module.GetGithubUser(tokenResponse.AccessToken)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get github user: %v", err))
		return
	}

	user, _ := user_service.GetUser(ginCtx, githubUser.Login)
	if user == nil {
		user, err = user_service.CreateUser(ginCtx, user_model.User{
			Account:   githubUser.Login,
			Name:      githubUser.Name,
			Email:     &githubUser.Email,
			AvatarURL: githubUser.AvatarURL,
		})
		if err != nil {
			gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to create user: %v", err))
			return
		}
	}

	ls, err := user_service.StartLoginSession(ginCtx, user.Account)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to start login session: %v", err))
		return
	}

	middleware.SetLoginSessionKeyCookie(ginCtx, ls.Key)

	ginCtx.JSON(200, nil)
}

func loginGithub(ginCtx *gin.Context) {
	u, err := auth_module.GetGithubOauthEntryURL(callbackURL)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get github oauth entry url: %v", err))
		return
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
		gin_utils.NewInvalidParamError(ginCtx, "body", "invalid body")
		return
	}

	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUserByAccountPassword(db, body.Account, body.Password)
	if err != nil {
		gin_utils.NewUnauthorizedError(ginCtx, "account or password incorrect")
		return
	}

	ls, err := user_service.StartLoginSession(ginCtx, user.Account)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to login: %v", err))
		return
	}
	middleware.SetLoginSessionKeyCookie(ginCtx, ls.Key)

	ginCtx.Status(http.StatusOK)
}
