package judger_agent

import (
	"io"
	"net/http"
	"net/url"
)

const judgerStateRoute = "/api/v1/state"

func (jc JudgerClient) GetState() (string, error) {
	url, err := url.JoinPath(jc.Host, judgerStateRoute)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	innerRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	innerResponse, err := client.Do(innerRequest)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(innerResponse.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
