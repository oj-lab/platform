package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/platform/cmd/web_server/handler"
	"github.com/oj-lab/platform/cmd/web_server/middleware"
	sloggin "github.com/samber/slog-gin"

	core_module "github.com/oj-lab/platform/modules/core"
)

const (
	servicePortConfigKey  = "service.port"
	swaggerOnConfigKey    = "service.swagger_on"
	frontendDistConfigKey = "service.frontend_dist"
)

var (
	servicePort  uint
	swaggerOn    bool
	frontendDist string
)

func init() {
	servicePort = core_module.Config.GetUint(servicePortConfigKey)
	swaggerOn = core_module.Config.GetBool(swaggerOnConfigKey)
	frontendDist = core_module.Config.GetString(frontendDistConfigKey)
}

func GetProjectDir() string {
	_, b, _, _ := runtime.Caller(0)
	projectDir := filepath.Join(filepath.Dir(b), "..", "..")

	return projectDir
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(sloggin.New(slog.Default().With("module", "gin")))
	r.Use(gin.Recovery())
	r.Use(middleware.HandleError)

	baseRouter := r.Group("/")
	if frontendDist != "" {
		// If dist folder is not empty, serve frontend
		if _, err := os.Stat(frontendDist); os.IsNotExist(err) {
			slog.Warn("Frontend dist is set but folder not found")
		} else {
			slog.Info("Serving frontend...")
			r.LoadHTMLFiles(frontendDist + "/index.html")
			handler.SetupFrontendRoute(baseRouter, frontendDist)
			r.NoRoute(handler.RenderHTML)
		}
	}

	if swaggerOn {
		slog.Info("Serving swagger Doc...")
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
