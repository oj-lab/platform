package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/service/problem"
	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/v1/user")
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
		g.GET("/:slug", problem.GetProblemInfo)
		g.PUT("/:slug/package", problem.PutProblemPackage)
		g.POST("/:slug/judge", problem.Judge)
	}
}
