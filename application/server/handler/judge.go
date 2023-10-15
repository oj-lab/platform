package handler

import (
	"github.com/OJ-lab/oj-lab-services/service"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/gin-gonic/gin"
)

func SetupJudgeRoute(baseRoute *gin.RouterGroup) {
	g := baseRoute.Group("/judge")
	{
		g.POST("/add-judger", postJudger)
	}
}

func postJudger(ginCtx *gin.Context) {
	judger := model.Judger{}
	if err := ginCtx.ShouldBindJSON(&judger); err != nil {
		ginCtx.Error(err)
		return
	}

	if err := service.AddJudger(ginCtx, judger); err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.JSON(200, gin.H{
		"message": "success",
	})
}
