package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupFrontendRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/")
	{
		// Page routing
		g.GET("/", RenderHTML)
		g.GET("/problem", RenderHTML)
		g.GET("/problem/:slug", RenderHTML)
		g.GET("/admin", RenderHTML)
		g.GET("/admin/problem", RenderHTML)

		g.Static("manifest.json", "frontend/dist/manifest.json")
		g.Static("/assets", "frontend/dist/assets")
		g.Static("/avatars", "frontend/dist/avatars")
		g.Static("/images", "frontend/dist/images")
		g.Static("/katex-dist", "frontend/dist/katex-dist")
	}
}

func RenderHTML(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
