package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/platform/cmd/web_server/handler"
	"github.com/oj-lab/platform/cmd/web_server/middleware"

	config_module "github.com/oj-lab/platform/modules/config"

	log_module "github.com/oj-lab/platform/modules/log"
)

const (
	serviceForceConsoleColorProp = "service.force_console_color"
	servicePortProp              = "service.port"
	serviceModeProp              = "service.mode"
	swaggerOnProp                = "service.swagger_on"
	frontendDistProp             = "service.frontend_dist"
)

var (
	serviceForceConsoleColor bool
	servicePort              uint
	serviceMode              string
	swaggerOn                bool
	frontendDist             string
)

func init() {
	serviceForceConsoleColor = config_module.AppConfig().GetBool(serviceForceConsoleColorProp)
	servicePort = config_module.AppConfig().GetUint(servicePortProp)
	serviceMode = config_module.AppConfig().GetString(serviceModeProp)
	swaggerOn = config_module.AppConfig().GetBool(swaggerOnProp)
	frontendDist = config_module.AppConfig().GetString(frontendDistProp)
}

func GetProjectDir() string {
	_, b, _, _ := runtime.Caller(0)
	projectDir := filepath.Join(filepath.Dir(b), "..", "..")

	return projectDir
}

func main() {
	if serviceForceConsoleColor {
		gin.ForceConsoleColor()
	}
	r := gin.Default()
	r.Use(middleware.HandleError)
	gin.SetMode(serviceMode)

	baseRouter := r.Group("/")
	if frontendDist != "" {
		// If dist folder is not empty, serve frontend
		if _, err := os.Stat(frontendDist); os.IsNotExist(err) {
			log_module.AppLogger().Warn("Frontend dist is set but folder not found")
		} else {
			log_module.AppLogger().Info("Serving frontend...")
			r.LoadHTMLFiles(frontendDist + "/index.html")
			handler.SetupFrontendRoute(baseRouter, frontendDist)
			r.NoRoute(handler.RenderHTML)
		}
	}

	if swaggerOn {
		log_module.AppLogger().Info("Serving swagger Doc...")
		handler.SetupSwaggoRouter(baseRouter)
	}
	handler.SetupAuthRouter(baseRouter)

	apiRouter := r.Group("/api/v1")
	handler.SetupUserRouter(apiRouter)
	handler.SetupProblemRouter(apiRouter)
	handler.SetupEventRouter(apiRouter)
	handler.SetupJudgeRouter(apiRouter)
	handler.SetupJudgeTaskRouter(apiRouter)
	handler.SetupJudgeResultRouter(apiRouter)
	handler.SetupRankRouter(apiRouter)

	err := r.Run(fmt.Sprintf(":%d", servicePort))
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
