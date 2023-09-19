package handler

import (
	swaggoGen "github.com/OJ-lab/oj-lab-services/application/server/swaggo-gen"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	servicePortProp = "service.port"
)

var (
	swaggerHost string
)

func SetupSwaggoRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func init() {
	sevicePort := viper.GetString(servicePortProp)
	swaggerHost = "localhost" + sevicePort
	println("Swagger host is set to: " + swaggerHost)
	// programmatically set swagger info
	swaggoGen.SwaggerInfo.Title = "OJ Lab Services API"
	swaggoGen.SwaggerInfo.Version = "1.0"
	swaggoGen.SwaggerInfo.Host = swaggerHost
	swaggoGen.SwaggerInfo.BasePath = "/api/v1"
	swaggoGen.SwaggerInfo.Schemes = []string{"http"}
}
