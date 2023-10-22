package main

import (
	"path/filepath"
	"runtime"

	"github.com/OJ-lab/oj-lab-services/application/server/handler"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/core/middleware"
	"github.com/gin-gonic/gin"
)

const (
	servicePortProp   = "service.port"
	serviceModeProp   = "service.mode"
	swaggerOnProp     = "service.swagger_on"
	serveFrontendProp = "service.serve_front"
)

var (
	servicePort   string
	serviceMode   string
	swaggerOn     bool
	serveFrontend bool
)

func init() {
	servicePort = core.AppConfig.GetString(servicePortProp)
	serviceMode = core.AppConfig.GetString(serviceModeProp)
	swaggerOn = core.AppConfig.GetBool(swaggerOnProp)
	serveFrontend = core.AppConfig.GetBool(serveFrontendProp)
}

func GetProjectDir() string {
	_, b, _, _ := runtime.Caller(0)
	projectDir := filepath.Join(filepath.Dir(b), "..", "..")

	return projectDir
}

func main() {
	r := gin.Default()
	r.Use(middleware.HandleError)
	gin.SetMode(serviceMode)

	baseRouter := r.Group("/")
	if serveFrontend {
		core.AppLogger().Info("Serving frontend...")
		r.LoadHTMLFiles("./frontend/dist/index.html")
		handler.SetupFrontendRoute(baseRouter)
	}

	if swaggerOn {
		core.AppLogger().Info("Serving swagger Doc...")
		handler.SetupSwaggoRouter(baseRouter)
	}

	apiRouter := r.Group("/api/v1")
	handler.SetupUserRouter(apiRouter)
	handler.SetupProblemRoute(apiRouter)
	handler.SetupEventRouter(apiRouter)

	err := r.Run(servicePort)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
