package router

import (
	"github.com/OJ-lab/oj-lab-services/user-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupUserRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, this is user service")
	})

	r.POST("/login", service.Login)
}
