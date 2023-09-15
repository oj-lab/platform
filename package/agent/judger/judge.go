package judger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/OJ-lab/oj-lab-services/package/core"
)

const JUDGER_HOST_PROP = "judger.host"
const JUDGER_JUDGE_PATH = "/api/v1/judge"

var judgerHost string

func init() {
	judgerHost = core.AppConfig.GetString(JUDGER_HOST_PROP)
}

type JudgeRequest struct {
	Src         string `json:"src"`
	SrcLanguage string `json:"src_language"`
}

func PostJudgeSync(packagelug string, judgeRequest JudgeRequest) ([]map[string]interface{}, error) {
	url, err := url.JoinPath(judgerHost, JUDGER_JUDGE_PATH, packagelug)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := json.Marshal(judgeRequest)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	innerRequest, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	innerRequest.Header.Set("Content-Type", "core/json")
	innerRequest.Header.Set("Accept", "core/json")
	innerResponse, err := client.Do(innerRequest)
	if err != nil {
		return nil, err
	}
	defer innerResponse.Body.Close()

	innerBody := []map[string]interface{}{}
	if err := json.NewDecoder(innerResponse.Body).Decode(&innerBody); err != nil {
		return nil, err
	}

	return innerBody, nil
}
