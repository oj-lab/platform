package judge_model

import "github.com/oj-lab/platform/models"

type JudgeRank struct {
	Rank             int
	AvatarURL        string
	Name             string
	Points           int
	TotalSubmissions int
	AcceptRate       float32
}

var RankInfoSelection = append([]string{"user_account", "points", "total_submissions"}, models.MetaFieldsSelection...)
