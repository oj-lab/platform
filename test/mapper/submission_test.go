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

	fmt.Printf("%+v\n", submissionList)
}
