package mapper_test

import (
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
)

func TestJudgerMapper(t *testing.T) {
	db := gormAgent.GetDefaultDB()
	_, err := mapper.GetJudgerList(db)
	if err != nil {
		t.Error(err)
	}
}
