package main

import (
	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/user-service/router"
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
	servicePort = application.AppConfig.GetString(servicePortProp)
	serviceMode = application.AppConfig.GetString(serviceModeProp)
}

func main() {
	r := gin.Default()
	gin.SetMode(serviceMode)
	router.SetupUserRouter(r)

	err := r.Run(servicePort)
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
