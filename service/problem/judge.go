package problem

import (
	"github.com/gin-gonic/gin"

	"github.com/OJ-lab/oj-lab-services/packages/agent/judger"
)

func Judge(ctx *gin.Context) {
	packageSlug := ctx.Param("slug")
	judgeRequest := judger.JudgeRequest{}
	if err := ctx.ShouldBindJSON(&judgeRequest); err != nil {
		ctx.Error(err)
		return
	}

	body, err := judger.PostJudgeSync(packageSlug, judgeRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, body)
}
