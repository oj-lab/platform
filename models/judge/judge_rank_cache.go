package judge_model

import (
	"github.com/oj-lab/oj-lab-platform/models"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
)

// user contest summary rank info
type RankCache struct {
	models.MetaFields
	UserAccount      string          `json:"userAccount" gorm:"primaryKey"`
	User             user_model.User `json:"user"`
	Points           uint            `json:"points"`
	TotalSubmissions uint            `json:"totalSubmissions"`
}
