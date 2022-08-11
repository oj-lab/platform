package router

import (
	"github.com/OJ-lab/oj-lab-services/user-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupUserRouter(r *gin.Engine) {
	g := r.Group("/api/user")
	{
		g.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is user service")
		})

		g.GET("/:account", service.GetUserInfo)
		g.POST("/:account/login", service.Login)
		g.POST("/:account/delete", service.DeleteUser)

		g.GET("/current", service.GetCurrentUser)
		g.POST("/register", service.Register)

		g.GET("", service.FindUserInfos)
	}
}
