package handler

import (
	"github.com/OJ-lab/oj-lab-services/packages/mapper"
	"github.com/gin-gonic/gin"
)

func GetProblemInfo(ctx *gin.Context) {
	slug := ctx.Param("slug")

	problem, err := mapper.GetProblem(slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"slug":        problem.Slug,
		"title":       problem.Title,
		"description": problem.Description,
		"tags":        mapper.GetTagsList(problem),
	})
}
