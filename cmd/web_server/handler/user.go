package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/platform/cmd/web_server/middleware"
	user_model "github.com/oj-lab/platform/models/user"
	casbin_agent "github.com/oj-lab/platform/modules/agent/casbin"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
	user_service "github.com/oj-lab/platform/services/user"
)

func SetupUserRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/user")
	{
		g.GET("",
			middleware.HandleRequireLogin,
			middleware.BuildCasbinEnforceHandlerWithDomain("system"),
			getUserList,
		)
		g.GET("/current", middleware.HandleRequireLogin, getCurrentUser)
		g.POST("/:account/role",
			middleware.HandleRequireLogin,
			middleware.BuildCasbinEnforceHandlerWithDomain("system"),
			grantUserRole,
		)
		g.DELETE("/:account/role",
			middleware.HandleRequireLogin,
			middleware.BuildCasbinEnforceHandlerWithDomain("system"),
			revokeUserRole,
		)
		g.POST("/logout", logout)
	}
}

func AddUserCasbinPolicies() error {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	_, err := enforcer.AddPolicies([][]string{
		{
			casbin_agent.RoleSubjectPrefix + `admin`, `true`, `system`,
			`/api/v1/user/*any`, strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, "|"),
			"allow",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserList(ginCtx *gin.Context) {
	limit, err := gin_utils.QueryInt(ginCtx, "limit", 10)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "limit", err.Error())
		return
	}
	offset, err := gin_utils.QueryInt(ginCtx, "offset", 0)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "offset", err.Error())
		return
	}

	users, total, err := user_service.GetUserList(ginCtx, user_model.GetUserOptions{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get user list: %v", err))
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
func getCurrentUser(ginCtx *gin.Context) {
	ls, err := middleware.GetLoginSessionFromGinCtx(ginCtx)
	if err != nil {
		gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
		return
	}
	user, err := user_service.GetUser(ginCtx, ls.Key.Account)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to get user: %v", err))
		return
	}

	ginCtx.JSON(http.StatusOK, user)
}

type grantUserRoleBody struct {
	Role   string `json:"role" example:"admin"`
	Domain string `json:"domain" example:"system"`
}

func grantUserRole(ginCtx *gin.Context) {
	account := ginCtx.Param("account")
	body := &grantUserRoleBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "body", "invalid body")
		return
	}

	err = user_service.GrantUserRole(ginCtx, account, body.Role, body.Domain)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to grant user role: %v", err))
		return
	}

	ginCtx.Status(http.StatusOK)
}

type revokeUserRoleBody struct {
	Role   string `json:"role" example:"admin"`
	Domain string `json:"domain" example:"system"`
}

func revokeUserRole(ginCtx *gin.Context) {
	account := ginCtx.Param("account")
	body := &revokeUserRoleBody{}
	err := ginCtx.BindJSON(body)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "body", "invalid body")
		return
	}

	err = user_service.RevokeUserRole(ginCtx, account, body.Role, body.Domain)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to revoke user role: %v", err))
		return
	}

	ginCtx.Status(http.StatusOK)
}

func logout(ginCtx *gin.Context) {
	ls, err := middleware.GetLoginSessionFromGinCtx(ginCtx)
	if err != nil {
		gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
		return
	}

	err = ls.DelInRedis(ginCtx)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, fmt.Sprintf("failed to Del session in redis: %v", err))
		return
	}

	ginCtx.Status(http.StatusOK)
}
