package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupFrontendRoute(baseRoute *gin.RouterGroup, frontendDist string) {
	g := baseRoute.Group("/")
	{
		// Page routing
		g.GET("/", RenderHTML)
		g.GET("/problem", RenderHTML)
		g.GET("/problem/:slug", RenderHTML)
		g.GET("/admin", RenderHTML)
		g.GET("/admin/problem", RenderHTML)

		g.Static("manifest.json", fmt.Sprintf("%s/manifest.json", frontendDist))
		g.Static("/assets", fmt.Sprintf("%s/assets", frontendDist))
		g.Static("/avatars", fmt.Sprintf("%s/avatars", frontendDist))
		g.Static("/images", fmt.Sprintf("%s/images", frontendDist))
		g.Static("/katex-dist", fmt.Sprintf("%s/katex-dist", frontendDist))
	}
}

func RenderHTML(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
