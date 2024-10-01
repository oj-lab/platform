package judge_service

import (
	"context"

	judge_model "github.com/oj-lab/platform/models/judge"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
)

func GetRankList(
	_ context.Context,
	account *string,
	limit, offset *int,
) ([]judge_model.JudgeRank, int64, error) {
	db := gorm_agent.GetDefaultDB()
	getOptions := judge_model.GetRankCacheOptions{
		// Selection: judge_model.RankCacheInfoSelection,
		Limit:  limit,
		Offset: offset,
	}

	total, err := judge_model.CountRankByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	rankCacheList, err := judge_model.GetRankCacheListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	var rankList []judge_model.JudgeRank

	for i, rankCache := range rankCacheList {
		if rankCache.TotalSubmissions == 0 {
			continue
		}
		rankList = append(rankList, judge_model.JudgeRank{
			Rank:             i + *offset + 1,
			User:             rankCache.User,
			Points:           rankCache.Points,
			TotalSubmissions: rankCache.TotalSubmissions,
			AcceptRate:       float32(rankCache.Points) / float32(rankCache.TotalSubmissions),
		})
	}
	return rankList, total, nil
}
