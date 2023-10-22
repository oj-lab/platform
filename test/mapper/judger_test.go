package mapper_test

import (
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
)

func TestJudgerMapper(t *testing.T) {
	db := gormAgent.GetDefaultDB()
	_, err := mapper.GetJudgerList(db)
	if err != nil {
		t.Error(err)
	}
}
