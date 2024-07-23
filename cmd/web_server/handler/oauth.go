package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/modules"
	auth_module "github.com/oj-lab/oj-lab-platform/modules/auth"
	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
)

func SetupOauthRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/oauth")
	{
		g.GET("/github/callback", githubCallback)
		g.Any("/github", loginGithub)
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
