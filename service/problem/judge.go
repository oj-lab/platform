package problem

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/gin-gonic/gin"
)

const JUDGER_HOST_PROP = "judger.host"
const JUDGER_JUDGE_PATH = "/api/v1/judge"

var judgerHost string

func init() {
	judgerHost = application.AppConfig.GetString(JUDGER_HOST_PROP)
}

type JudgeRequest struct {
	Src         string `json:"src"`
	SrcLanguage string `json:"src_language"`
}

func Judge(ctx *gin.Context) {
	packageSlug := ctx.Param("slug")

	judgeRequest := JudgeRequest{}
	if err := ctx.ShouldBindJSON(&judgeRequest); err != nil {
		ctx.Error(err)
		return
	}

	url, err := url.JoinPath(judgerHost, JUDGER_JUDGE_PATH, packageSlug)
	if err != nil {
		ctx.Error(err)
		return
	}
	payloadBytes, err := json.Marshal(judgeRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	client := &http.Client{}
	innerRequest, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		ctx.Error(err)
		return
	}
	innerRequest.Header.Set("Content-Type", "application/json")
	innerRequest.Header.Set("Accept", "application/json")
	innerResponse, err := client.Do(innerRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	defer innerResponse.Body.Close()

	innerBody := []map[string]interface{}{}
	if err := json.NewDecoder(innerResponse.Body).Decode(&innerBody); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(innerResponse.StatusCode, innerBody)
}
