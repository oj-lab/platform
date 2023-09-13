package main

import (
	"github.com/OJ-lab/oj-lab-services/application/server/handler"
	"github.com/OJ-lab/oj-lab-services/packages/core"
	"github.com/gin-gonic/gin"
)

const (
	servicePortProp = "service.port"
	serviceModeProp = "service.mode"
)

var (
	servicePort string
	serviceMode string
)

func init() {
	servicePort = core.AppConfig.GetString(servicePortProp)
	serviceMode = core.AppConfig.GetString(serviceModeProp)
}

func main() {
	r := gin.Default()
	r.Use(core.HandleError)
	gin.SetMode(serviceMode)
	handler.SetupUserRouter(r)
	handler.SetupProblemRoute(r)

	err := r.Run(servicePort)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
