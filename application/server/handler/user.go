package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/core/middleware"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/v1/user")
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
	Account  string `json:"account"`
	Password string `json:"password"`
}

func login(ginCtx *gin.Context) {
	body := &loginBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		core.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
		return
	}

	lsId, svcErr := service.StartLoginSession(ginCtx, body.Account, body.Password)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}
	middleware.SetLoginSessionCookie(ginCtx, *lsId)

	ginCtx.String(http.StatusOK, "")
}

func me(ginCtx *gin.Context) {
	ls := middleware.GetLoginSession(ginCtx)
	if ls == nil {
		core.NewUnauthorizedError("not logined").AppendToGin(ginCtx)
		return
	}
	user, svcErr := service.GetUser(ginCtx, ls.Account)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(http.StatusOK, user)
}

func checkUserExist(ginCtx *gin.Context) {
	account := ginCtx.Query("account")
	if account == "" {
		core.NewInvalidParamError("account", "account cannot be empty").AppendToGin(ginCtx)
		return
	}

	exist, err := service.CheckUserExist(ginCtx, account)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{
		"exist": exist,
	})
}
