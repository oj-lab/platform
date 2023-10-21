package mapper_test

import (
	"fmt"
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func TestSubmissionMapper(t *testing.T) {
	db := gormAgent.GetDefaultDB()
	getOptions := mapper.GetSubmissionOptions{
		UserAccount: func() *string { s := "admin"; return &s }(),
		OrderByColumns: []model.OrderByColumnOption{{
			Column: string(model.JudgeTaskSubmissionSortByColumnCreateAt),
			Desc:   true,
		}},
	}
	submissionList, _, err := mapper.GetSubmissionListByOptions(db, getOptions)
	if err != nil {
		t.Error(err)
	}
	if len(submissionList) < 2 {
		t.Error("submission list should not be less than 2")
	}
	if submissionList[0].MetaFields.CreateAt.Before(*submissionList[1].MetaFields.CreateAt) {
		t.Error("submission list should be sorted by create_at desc")
	}

	fmt.Printf("%+v\n", submissionList)
}
