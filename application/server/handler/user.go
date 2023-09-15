package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/package/core"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/v1/user")
	{
		g.GET("/health", func(ginCtx *gin.Context) {
			ginCtx.String(http.StatusOK, "Hello, this is user service")
		})
		g.GET("/me", func(ginCtx *gin.Context) {
			ginCtx.String(http.StatusOK, "WIP")
		})
		g.GET("/check-exist", checkUserExist)
	}
}

func checkUserExist(ctx *gin.Context) {
	account := ctx.Query("account")
	if account == "" {
		core.NewInvalidParamError("account", "account cannot be empty").AppendToGin(ctx)
		return
	}

	exist, err := service.CheckUserExist(account)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"exist": exist,
	})
}
