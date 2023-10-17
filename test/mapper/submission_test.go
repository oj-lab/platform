package mapper_test

import (
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func TestSubmissionMapper(t *testing.T) {
	getOptions := mapper.GetSubmissionOptions{
		UserAccount: func() *string { s := "admin"; return &s }(),
		OrderByColumns: []model.OrderByColumnOption{{
			Column: string(model.JudgeTaskSubmissionSortByColumnCreateAt),
			Desc:   true,
		}},
	}
	submissionList, _, err := mapper.GetSubmissionListByOptions(getOptions)
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
