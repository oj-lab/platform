package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/core/agent/judger"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/gin-gonic/gin"
)

func SetupProblemRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/problem")
	{
		g.GET("/greet", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is problem service")
		})
		g.GET("/:slug", getProblemInfo)
		g.PUT("/:slug/package", putProblemPackage)
		g.POST("/:slug/judge", judge)
	}
}

func getProblemInfo(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")

	problemInfo, err := service.GetProblemInfo(ginCtx, slug)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"slug":        problemInfo.Slug,
		"title":       problemInfo.Title,
		"description": problemInfo.Description,
		"tags":        mapper.GetTagsList(*problemInfo),
	})
}

func putProblemPackage(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	file, err := ginCtx.FormFile("file")
	if err != nil {
		ginCtx.Error(err)
		return
	}
	zipFile := "/tmp/" + slug + ".zip"
	if err := ginCtx.SaveUploadedFile(file, zipFile); err != nil {
		ginCtx.Error(err)
		return
	}

	service.PutProblemPackage(ginCtx, slug, zipFile)

	ginCtx.Done()
}

func judge(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	judgeRequest := judger.JudgeRequest{}
	if err := ginCtx.ShouldBindJSON(&judgeRequest); err != nil {
		ginCtx.Error(err)
		return
	}

	body, err := service.Judge(ginCtx, slug, judgeRequest)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, body)
}
