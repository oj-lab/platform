package main

import (
	"os"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/OJ-lab/oj-lab-services/user-service/router"
	"github.com/OJ-lab/oj-lab-services/user-service/service"
	"github.com/gin-gonic/gin"
)

func main() {
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "config/ini/default.ini"
	}

	dataBaseSettings, err := config.GetDatabaseSettings(configPath)
	if err != nil {
		panic("failed to get database settings")
	}
	jwtSettings, err := config.GetJWTSettings(configPath)
	if err != nil {
		panic("failed to get jwt settings")
	}
	serviceSettings, err := config.GetServiceSettings(configPath)
	if err != nil {
		panic("failed to get service settings")
	}

	business.OpenDBConnection(*dataBaseSettings)
	business.SetupJWTSettings(jwtSettings)
	service.SetupServiceSetting(serviceSettings)

	r := gin.Default()
	gin.SetMode(serviceSettings.Mode)
	router.SetupUserRouter(r)

	err = r.Run(serviceSettings.Port)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
