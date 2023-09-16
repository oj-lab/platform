package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/core/agent/judger"
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/gin-gonic/gin"
)

func SetupProblemRoute(r *gin.Engine) {
	g := r.Group("/api/v1/problem")
	{
		g.GET("/greet", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is problem service")
		})
		g.GET("/:slug", getProblemInfo)
		g.PUT("/:slug/package", putProblemPackage)
		g.POST("/:slug/judge", judge)
	}
}

func getProblemInfo(ctx *gin.Context) {
	slug := ctx.Param("slug")

	problemInfo, err := service.GetProblemInfo(slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"slug":        problemInfo.Slug,
		"title":       problemInfo.Title,
		"description": problemInfo.Description,
		"tags":        mapper.GetTagsList(*problemInfo),
	})
}

func putProblemPackage(ctx *gin.Context) {
	slug := ctx.Param("slug")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}
	zipFile := "/tmp/" + slug + ".zip"
	if err := ctx.SaveUploadedFile(file, zipFile); err != nil {
		ctx.Error(err)
		return
	}

	service.PutProblemPackage(slug, zipFile)

	ctx.Done()
}

func judge(ctx *gin.Context) {
	slug := ctx.Param("slug")
	judgeRequest := judger.JudgeRequest{}
	if err := ctx.ShouldBindJSON(&judgeRequest); err != nil {
		ctx.Error(err)
		return
	}

	body, err := service.Judge(slug, judgeRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, body)
}
