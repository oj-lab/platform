package judge_model

import (
	"github.com/oj-lab/platform/models"
	user_model "github.com/oj-lab/platform/models/user"
)

// user contest summary rank info
type JudgeRankCache struct {
	models.MetaFields
	UserAccount      string          `json:"userAccount" gorm:"primaryKey"`
	User             user_model.User `json:"user"`
	Points           int             `json:"points"`
	TotalSubmissions int             `json:"totalSubmissions"`
}

func NewJudgeRankCache(userAccount string) JudgeRankCache {
	return JudgeRankCache{
		UserAccount:      userAccount,
		Points:           0,
		TotalSubmissions: 0,
	}
}

var RankCacheInfoSelection = append([]string{"points", "total_submissions"}, models.MetaFieldsSelection...)
