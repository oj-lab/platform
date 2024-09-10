package judge_model

import user_model "github.com/oj-lab/platform/models/user"

type JudgeRank struct {
	Rank             int             `json:"rank"`
	User             user_model.User `json:"user"`
	Points           int             `json:"acceptCount"`
	TotalSubmissions int             `json:"submitCount"`
	AcceptRate       float32         `json:"acceptRate"`
}
