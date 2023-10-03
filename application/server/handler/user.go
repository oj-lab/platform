package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/core/middleware"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		g.GET("/stream", Stream)
	}
}

// Stream
//
//	@Summary		Stream
//	@Description	Stream
//	@Tags			user
//	@Router			/user/stream [get]
//	@Accept			text/event-stream
//	@Produce		text/event-stream
//	@Success		200	{string}	string	"data: {message}"
//	@Router			/user/stream [get]
func Stream(ginCtx *gin.Context) {
	ginCtx.Header("Content-Type", "text/event-stream")
	ginCtx.Header("Cache-Control", "no-cache")

	ginCtx.Stream(func(w io.Writer) bool {
		// 每秒钟向客户端发送一条消息
		logrus.Info("send message")
		fmt.Fprintf(w, "data: %s\n\n", time.Now().String())
		time.Sleep(1 * time.Second)
		return true
	})
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
		core.NewInvalidParamError("body", "invalid body").AppendToGin(ginCtx)
		return
	}

	lsId, svcErr := service.StartLoginSession(ginCtx, body.Account, body.Password)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}
	middleware.SetLoginSessionCookie(ginCtx, *lsId)

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
