package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/user")
	{
		g.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is user service")
		})
	}
}

func SetupProblemRoute(r *gin.Engine) {
	g := r.Group("/api/v1/problem")
	{
		g.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is problem service")
		})
		g.GET("")
	}
}
