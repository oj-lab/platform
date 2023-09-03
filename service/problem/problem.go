package problem

import (
	"github.com/OJ-lab/oj-lab-services/packages/mapper"
	"github.com/OJ-lab/oj-lab-services/packages/utils"
	"github.com/gin-gonic/gin"
)

func GetProblemInfo(c *gin.Context) {
	slug := c.Param("slug")
	problem, err := mapper.GetProblem(slug)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"slug":        problem.Slug,
		"title":       problem.Title,
		"description": problem.Description,
		"tags":        mapper.GetTagsList(problem),
	})
}
