package core_test

import (
	"testing"

	judgerAgent "github.com/OJ-lab/oj-lab-services/core/agent/judger"
)

func TestJudgerGetState(t *testing.T) {
	judgerClient := judgerAgent.JudgerClient{
		Host: "http://localhost:8000",
	}

	state, err := judgerClient.GetState()
	if err != nil {
		t.Error(err)
	}
	println(state)
}
