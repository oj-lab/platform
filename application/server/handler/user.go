package handler

import (
	"net/http"

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
