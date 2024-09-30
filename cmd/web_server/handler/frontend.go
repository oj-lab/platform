package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/platform/cmd/web_server/middleware"
	casbin_agent "github.com/oj-lab/platform/modules/agent/casbin"
)

func SetupFrontendRoute(baseRoute *gin.RouterGroup, frontendDist string) {
	g := baseRoute.Group("/")
	{
		// Page routing
		g.GET("/", RenderHTML)
		g.GET("/problems", RenderHTML)
		g.GET("/problems/:slug", RenderHTML)
		g.GET("/judges", RenderHTML)
		g.GET("/judges/:uid", RenderHTML)
		g.GET("/rank", RenderHTML)
		g.GET("/admin/*any",
			middleware.HandleRequireLogin,
			middleware.BuildCasbinEnforceHandlerWithDomain("system"),
			RenderHTML,
		)

		// Static file routing
		g.StaticFile("manifest.json", fmt.Sprintf("%s/manifest.json", frontendDist))
		g.Static("/assets", fmt.Sprintf("%s/assets", frontendDist))
		g.Static("/avatars", fmt.Sprintf("%s/avatars", frontendDist))
		g.Static("/images", fmt.Sprintf("%s/images", frontendDist))
		g.Static("/katex-dist", fmt.Sprintf("%s/katex-dist", frontendDist))
	}
}

func AddFrontendPagePolicies() error {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	_, err := enforcer.AddPolicies([][]string{
		{
			casbin_agent.RoleSubjectPrefix + `admin`, `true`, `system`,
			`/admin/*any`, strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, "|"),
			"allow",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func RenderHTML(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
