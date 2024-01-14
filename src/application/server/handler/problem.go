package handler

import (
	"net/http"

	"github.com/OJ-lab/oj-lab-services/src/service"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
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
		g.POST("/:slug/submission", postSubmission)
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

type PostSubmissionBody struct {
	Code     string                   `json:"code" binding:"required"`
	Language model.SubmissionLanguage `json:"language" binding:"required"`
}

// postSubmission
//
//	@Router			/problem/{slug}/submission [post]
//	@Summary		Post submission
//	@Description	Post submission
//	@Tags			problem
//	@Accept			json
//	@Param			slug			path	string				true	"problem slug"
//	@Param			judgeRequest	body	PostSubmissionBody	true	"judge request"
func postSubmission(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	body := PostSubmissionBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		ginCtx.Error(err)
		return
	}

	submission := model.NewSubmission("", slug, body.Code, body.Language)
	result, svcErr := service.CreateJudgeTaskSubmission(ginCtx, submission)
	if svcErr != nil {
		svcErr.AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, result)
}
