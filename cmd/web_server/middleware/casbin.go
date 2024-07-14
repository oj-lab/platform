package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/modules"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

func BuildCasbinEnforceHandlerWithDomain(domain string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		enforcer := casbin_agent.GetDefaultCasbinEnforcer()
		path := ginCtx.Request.URL.Path
		method := ginCtx.Request.Method
		ls, err := GetLoginSessionFromGinCtx(ginCtx)
		if err != nil {
			modules.NewUnauthorizedError("cannot load login session from cookie").AppendToGin(ginCtx)
			return
		}

		allow, err := enforcer.Enforce(ls.Key.Account, "_", domain, path, method)
		if err != nil {
			log.AppLogger().Errorf("Failed to enforce: %v", err)
			modules.NewInternalError("Failed to enforce").AppendToGin(ginCtx)
			return
		}
		if !allow {
			modules.NewUnauthorizedError("Unauthorized").AppendToGin(ginCtx)
			return
		}
		ginCtx.Next()
	}
}
