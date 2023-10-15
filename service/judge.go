package service

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/service/business"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func AddJudger(ctx context.Context, judger model.Judger) error {
	// TODO: validate judger

	return business.EnqueueWaitJudgerTask(ctx, judger)
}
