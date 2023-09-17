package main

import (
	"github.com/OJ-lab/oj-lab-services/application/server/handler"
	"github.com/sirupsen/logrus"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/core/middleware"
	"github.com/gin-gonic/gin"
)

const (
	servicePortProp  = "service.port"
	serviceModeProp  = "service.mode"
	serveSwaggerProp = "service.swagger_on"
)

var (
	servicePort string
	serviceMode string
	swaggerOn   bool
)

func init() {
	servicePort = core.AppConfig.GetString(servicePortProp)
	serviceMode = core.AppConfig.GetString(serviceModeProp)
	swaggerOn = core.AppConfig.GetBool(serveSwaggerProp)
}

func main() {
	r := gin.Default()
	r.Use(middleware.HandleError)
	gin.SetMode(serviceMode)

	baseRouter := r.Group("/")
	if swaggerOn {
		logrus.Info("Serving swagger Doc")
		handler.SetupSwaggoRouter(baseRouter)
	}

	apiRouter := r.Group("/api/v1")
	handler.SetupUserRouter(apiRouter)
	handler.SetupProblemRoute(apiRouter)

	err := r.Run(servicePort)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
