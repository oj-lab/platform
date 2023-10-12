package mapper_test

import (
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/service/mapper"
)

func TestSubmissionMapper(t *testing.T) {
	getOptions := mapper.GetSubmissionOptions{
		UserAccount: func() *string { s := "admin"; return &s }(),
	}
	submissionList, count, err := mapper.GetSubmissionListByOptions(getOptions)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("submission count should be 1")
	}
	if len(submissionList) != 1 {
		t.Error("submission list length should be 1")
	}
	fmt.Printf("%+v\n", submissionList[0])
}
