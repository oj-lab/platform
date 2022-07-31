package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupUserRouter(r *gin.Engine) {
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		if lastname == "" {
			c.String(http.StatusBadRequest, "lastname is required")
		}

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
}
