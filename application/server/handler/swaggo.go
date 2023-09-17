package handler

import (
	swaggoGen "github.com/OJ-lab/oj-lab-services/application/server/swaggo-gen"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupSwaggoRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func init() {
	// programmatically set swagger info
	swaggoGen.SwaggerInfo.Title = "OJ Lab Services API"
	swaggoGen.SwaggerInfo.Version = "1.0"
	swaggoGen.SwaggerInfo.Host = "localhost:8080"
	swaggoGen.SwaggerInfo.BasePath = "/api/v1"
	swaggoGen.SwaggerInfo.Schemes = []string{"http"}
}
