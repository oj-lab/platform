package main

import (
	"path/filepath"
	"runtime"

	"github.com/OJ-lab/oj-lab-services/src/application/server/handler"
	"github.com/OJ-lab/oj-lab-services/src/core"

	"github.com/OJ-lab/oj-lab-services/src/core/middleware"
	"github.com/gin-gonic/gin"
)

const (
	servicePortProp  = "service.port"
	serviceModeProp  = "service.mode"
	swaggerOnProp    = "service.swagger_on"
	frontendDistProp = "service.frontend_dist"
)

var (
	servicePort  string
	serviceMode  string
	swaggerOn    bool
	frontendDist string
)

func init() {
	servicePort = core.AppConfig.GetString(servicePortProp)
	serviceMode = core.AppConfig.GetString(serviceModeProp)
	swaggerOn = core.AppConfig.GetBool(swaggerOnProp)
	frontendDist = core.AppConfig.GetString(frontendDistProp)
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
	if frontendDist != "" {
		core.AppLogger().Info("Serving frontend...")
		r.LoadHTMLFiles(frontendDist + "/index.html")
		handler.SetupFrontendRoute(baseRouter, frontendDist)
	}

	if swaggerOn {
		core.AppLogger().Info("Serving swagger Doc...")
		handler.SetupSwaggoRouter(baseRouter)
	}

	apiRouter := r.Group("/api/v1")
	handler.SetupUserRouter(apiRouter)
	handler.SetupProblemRoute(apiRouter)
	handler.SetupEventRouter(apiRouter)
	handler.SetupSubmissionRouter(apiRouter)
	handler.SetupJudgeRoute(apiRouter)

	err := r.Run(servicePort)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
