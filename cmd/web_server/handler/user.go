package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/cmd/web_server/middleware"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	"github.com/oj-lab/oj-lab-platform/modules"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
	user_service "github.com/oj-lab/oj-lab-platform/services/user"
)

func SetupUserRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/user")
	{
		g.GET("",
			middleware.HandleRequireLogin,
			middleware.BuildCasbinEnforceHandlerWithDomain("system"),
			GetUserList,
		)
		g.GET("/me", middleware.HandleRequireLogin, me)
	}
}

func AddUserCasbinPolicies() error {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	_, err := enforcer.AddPolicies([][]string{
		{`role_admin`, `true`, `system`, `/api/v1/user/*any`, http.MethodGet, "allow"},
	})
	if err != nil {
		return err
	}
	return nil
}

func GetUserList(ginCtx *gin.Context) {
	options := user_model.GetUserOptions{}
	limitStr := ginCtx.Query("limit")
	offsetStr := ginCtx.Query("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	options.Limit = func() *int { return &limit }()
	options.Offset = func() *int { return &offset }()

	users, total, err := user_service.GetUserList(ginCtx, options)
	if err != nil {
		modules.NewInternalError(fmt.Sprintf("failed to get user list: %v", err)).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
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
