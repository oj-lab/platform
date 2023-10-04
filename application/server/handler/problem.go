package handler

import (
	"net/http"

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
		g.GET("", getProblemInfoList)
		g.GET("/:slug", getProblem)
		g.PUT("/:slug/package", putProblemPackage)
		g.POST("/:slug/judge/task", postJudgeTask)
		g.POST("/:slug/judge", judge)
	}
}

func getProblem(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")

	problemInfo, err := service.GetProblem(ginCtx, slug)
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

// getProblemInfoList
//
//	@Router			/problem [get]
//	@Summary		Get problem list
//	@Description	Get problem list
//	@Tags			problem
//	@Accept			json
//	@Success		200
func getProblemInfoList(ginCtx *gin.Context) {
	problemInfoList, total, err := service.GetProblemInfoList(ginCtx)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"total": total,
		"list":  problemInfoList,
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

type judgeTaskBody struct {
	Src         string `json:"src"`
	SrcLanguage string `json:"src_language"`
}

// postJudgeTask
//
//	@Router			/problem/{slug}/judge/task [post]
//	@Summary		Post judge task
//	@Description	Post judge task
//	@Tags			problem
//	@Accept			json
//	@Param			slug			path	string			true	"problem slug"
//	@Param			judgeRequest	body	judgeTaskBody	true	"judge request"
func postJudgeTask(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	body := judgeTaskBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		ginCtx.Error(err)
		return
	}

	service.PostJudgeTask(ginCtx, slug, body.Src, body.SrcLanguage)

	ginCtx.Done()
}

func judge(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	body := judgeTaskBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		ginCtx.Error(err)
		return
	}

	responseBody, err := service.Judge(ginCtx, slug, body.Src, body.SrcLanguage)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, responseBody)
}
