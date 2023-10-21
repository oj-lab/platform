package mapper_test

import (
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
)

func TestJudgerMapper(t *testing.T) {
	db := gormAgent.GetDefaultDB()
	judgerList, err := mapper.GetJudgerList(db)
	if err != nil {
		t.Error(err)
	}

	if len(judgerList) != 1 {
		t.Error("judger list should be 1")
	}
}
