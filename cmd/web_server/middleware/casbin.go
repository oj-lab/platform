package middleware

import (
	"github.com/gin-gonic/gin"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
	gin_utils "github.com/oj-lab/oj-lab-platform/modules/utils/gin"
)

func BuildCasbinEnforceHandlerWithDomain(domain string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		enforcer := casbin_agent.GetDefaultCasbinEnforcer()
		path := ginCtx.Request.URL.Path
		method := ginCtx.Request.Method
		ls, err := GetLoginSessionFromGinCtx(ginCtx)
		if err != nil {
			gin_utils.NewUnauthorizedError(ginCtx, "cannot load login session from cookie")
			ginCtx.Abort()
			return
		}

		allow, err := enforcer.Enforce(casbin_agent.UserSubjectPrefix+ls.Key.Account, "_", domain, path, method)
		if err != nil {
			log_module.AppLogger().Errorf("Failed to enforce: %v", err)
			gin_utils.NewInternalError(ginCtx, "Failed to enforce")
			ginCtx.Abort()
			return
		}
		if !allow {
			gin_utils.NewUnauthorizedError(ginCtx, "Unauthorized")
			ginCtx.Abort()
			return
		}
		ginCtx.Next()
	}
}
