package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/v1/user")
	{
		g.GET("/health", func(ginCtx *gin.Context) {
			ginCtx.String(http.StatusOK, "Hello, this is user service")
		})
		g.GET("/me", func(ginCtx *gin.Context) {
			ginCtx.String(http.StatusOK, "WIP")
		})
	}
}
