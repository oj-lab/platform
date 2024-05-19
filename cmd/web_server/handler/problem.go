package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	"github.com/oj-lab/oj-lab-platform/modules"
	judge_service "github.com/oj-lab/oj-lab-platform/services/judge"
	problem_service "github.com/oj-lab/oj-lab-platform/services/problem"
)

func SetupProblemRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/problem")
	{
		g.GET("/greet", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, this is problem service")
		})
		g.GET("", getProblemInfoList)
		g.PUT("", putProblem)
		g.GET("/:slug", getProblem)
		g.DELETE("/:slug", deleteProblem)
		g.GET("/:slug/check", checkProblemSlug)
		g.PUT("/:slug/package", putProblemPackage)
		g.POST("/:slug/submission", postSubmission)
	}
}

// getProblem
//
//	@Router			/problem/{slug} [get]
//	@Summary		Get a problem
//	@Description	Get a problem
//	@Tags			problem
//	@Accept			json
//	@Success		200
func getProblem(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")

	problemInfo, err := problem_service.GetProblem(ginCtx, slug)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"slug":        problemInfo.Slug,
		"title":       problemInfo.Title,
		"description": problemInfo.Description,
		"tags":        problem_model.GetTagsList(*problemInfo),
	})
}

// putProblem
//
//	@Router			/problem [put]
//	@Summary		Put a problem
//	@Description	Put a problem
//	@Tags			problem
//	@Accept			json
func putProblem(ginCtx *gin.Context) {
	problem := problem_model.Problem{}
	if err := ginCtx.ShouldBindJSON(&problem); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	err := problem_service.PutProblem(ginCtx, problem)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
}

// deleteProblem
//
//	@Router			/problem/{slug} [delete]
//	@Summary		Delete a problem
//	@Description	Delete a problem
//	@Tags			problem
//	@Accept			json
func deleteProblem(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")

	err := problem_service.DeleteProblem(ginCtx, slug)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
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
	problemInfoList, total, err := problem_service.GetProblemInfoList(ginCtx)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"total": total,
		"list":  problemInfoList,
	})
}

// putProblemPackage
//
//	@Router			/problem/{slug}/package [put]
//	@Summary		Put problem package
//	@Description	Put problem package
//	@Tags			problem
//	@Accept			json
//	@Param			slug	path	string	true	"problem slug"
func putProblemPackage(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	file, err := ginCtx.FormFile("file")
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	zipFile := "/tmp/" + slug + ".zip"
	if err := ginCtx.SaveUploadedFile(file, zipFile); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	err = problem_service.PutProblemPackage(ginCtx, slug, zipFile)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.Done()
}

// checkProblemSlug
//
//	@Router			/problem/{slug}/check [get]
//	@Summary		Check problem slug
//	@Description	Check problem slug
//	@Tags			problem
//	@Accept			json
//	@Success		200
//	@Param			slug	path	string	true	"problem slug"
func checkProblemSlug(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")

	valid, err := problem_service.CheckProblemSlug(ginCtx, slug)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"valid": valid,
	})
}

// PostSubmissionBody
//
//	@Description	The body of a submission request, containing the code and the language used for the submission.
//	@Property		code (string) required "The source code of the submission" minlength(1)
//	@Property		language (SubmissionLanguage) required "The programming language used for the submission"
type PostSubmissionBody struct {
	Code     string                         `json:"code" binding:"required"`
	Language judge_model.SubmissionLanguage `json:"language" binding:"required"`
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
		_ = ginCtx.Error(err)
		return
	}

	submission := judge_model.NewSubmission("", slug, body.Code, body.Language)
	result, err := judge_service.CreateJudgeTaskSubmission(ginCtx, submission)
	if err != nil {
		modules.NewInternalError(err.Error()).AppendToGin(ginCtx)
		return
	}

	ginCtx.JSON(200, result)
}
