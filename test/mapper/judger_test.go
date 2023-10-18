package mapper_test

import (
	"testing"

	"github.com/OJ-lab/oj-lab-services/service/mapper"
)

func TestJudgerMapper(t *testing.T) {
	judgerList, err := mapper.GetJudgerList()
	if err != nil {
		t.Error(err)
	}

	if len(judgerList) != 1 {
		t.Error("judger list should be 1")
	}
}
