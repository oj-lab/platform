package handler

import (
	"github.com/gin-gonic/gin"
	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	gin_utils "github.com/oj-lab/platform/modules/utils/gin"
	judge_service "github.com/oj-lab/platform/services/judge"
	problem_service "github.com/oj-lab/platform/services/problem"
)

func SetupProblemRouter(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/problem")
	{
		g.GET("", getProblemInfoList)
		g.PUT("", putProblem)
		g.GET("/:slug", getProblem)
		g.DELETE("/:slug", deleteProblem)
		g.GET("/:slug/check", checkProblemSlug)
		g.PUT("/:slug/package", putProblemPackage)
		g.POST("/:slug/judge", postJudge)
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

	problem, err := problem_service.GetProblem(ginCtx, slug)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, problem)
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
	limit, err := gin_utils.QueryInt(ginCtx, "limit", 10)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "limit", err.Error())
		return
	}
	offset, err := gin_utils.QueryInt(ginCtx, "offset", 0)
	if err != nil {
		gin_utils.NewInvalidParamError(ginCtx, "offset", err.Error())
		return
	}

	problemInfoList, total, err := problem_service.GetProblemInfoList(
		ginCtx,
		nil,
		&limit, &offset,
	)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, err.Error())
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

// PostJudgeBody
//
//	@Description	The body of a judge request, containing the code and the language used for the judge.
//	@Property		code (string) required "The source code of the judge" minlength(1)
//	@Property		language (ProgrammingLanguage) required "The programming language used for the judge"
type PostJudgeBody struct {
	Code     string                          `json:"code" binding:"required"`
	Language judge_model.ProgrammingLanguage `json:"language" binding:"required"`
}

// postJudge
//
//	@Router			/problem/{slug}/judge [post]
//	@Summary		Post judge
//	@Description	Post judge
//	@Tags			problem
//	@Accept			json
//	@Param			slug			path	string			true	"problem slug"
//	@Param			judgeRequest	body	PostJudgeBody	true	"judge request"
func postJudge(ginCtx *gin.Context) {
	slug := ginCtx.Param("slug")
	body := PostJudgeBody{}
	if err := ginCtx.ShouldBindJSON(&body); err != nil {
		_ = ginCtx.Error(err)
		return
	}

	judge := judge_model.NewJudge("", slug, body.Code, body.Language)
	result, err := judge_service.CreateJudge(ginCtx, judge)
	if err != nil {
		gin_utils.NewInternalError(ginCtx, err.Error())
		return
	}

	ginCtx.JSON(200, result)
}
