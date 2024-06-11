package judger_agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const judgerJudgeRoute = "/api/v1/judge"

type JudgeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

func (jc JudgerClient) PostJudgeSync(packagelug string, judgeRequest JudgeRequest) ([]map[string]interface{}, error) {
	url, err := url.JoinPath(jc.Host, judgerJudgeRoute, packagelug)
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
