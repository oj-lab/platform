package handler

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	log_module "github.com/oj-lab/platform/modules/log"
)

func SetupEventRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/event")
	{
		g.GET("", Stream)
	}
}

// Stream
//
//	@Summary		Stream
//	@Description	Stream
//	@Tags			user
//	@Router			/user/stream [get]
//	@Accept			text/event-stream
//	@Produce		text/event-stream
//	@Success		200	{string}	string	"data: {message}"
func Stream(ginCtx *gin.Context) {
	ginCtx.Header("Content-Type", "text/event-stream")
	ginCtx.Header("Cache-Control", "no-cache")

	counter := 0
	ginCtx.Stream(func(w io.Writer) bool {
		// With event type
		message := fmt.Sprintf("event: %s\ndata: %s\n\n", "eventType", time.Now().String())
		log_module.AppLogger().Infof("Send message:\n%s", message)
		fmt.Fprint(w, message)
		time.Sleep(1 * time.Second)
		counter++
		if counter > 5 {
			// Close connection (without event type)
			closeMessage := fmt.Sprintf("data: %s\n\n", "close")
			fmt.Fprint(w, closeMessage)
			return false
		}
		return true
	})
}
