package judge_service

import (
	"context"

	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func GetRankList(
	_ context.Context,
	account *string,
	limit, offset *int,
) ([]judge_model.JudgeRank, error) {
	db := gorm_agent.GetDefaultDB()
	getOptions := judge_model.GetRankOptions{
		Selection: judge_model.RankInfoSelection,
		Limit:     limit,
		Offset:    offset,
	}

	rankCacheList, err := judge_model.GetRankCacheListByOptions(db, getOptions)
	if err != nil {
		return nil, err
	}

	var rankList []judge_model.JudgeRank

	for i, rankCache := range rankCacheList {
		rankList = append(rankList, judge_model.JudgeRank{
			Rank:             i + *offset + 1,
			AvatarURL:        rankCache.User.AvatarURL,
			Name:             rankCache.User.Name,
			Points:           rankCache.Points,
			TotalSubmissions: rankCache.TotalSubmissions,
			AcceptRate:       float32(rankCache.Points) / float32(rankCache.TotalSubmissions),
		})
	}
	return rankList, nil
}
