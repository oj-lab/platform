package problem

import "github.com/oj-lab/oj-lab-platform/models"

type ProblemInfo struct {
	models.MetaFields
	ProblemSlug string        `json:"problem_slug" gorm:"primaryKey"`
	Title       string        `json:"title"`
	Tags        []*ProblemTag `json:"tags" gorm:"many2many:problem_info_tags;"`
}

type ProblemTag struct {
	models.MetaFields
	Name  string         `json:"name" gorm:"primaryKey"`
	Infos []*ProblemInfo `json:"problem_infos" gorm:"many2many:problem_info_tags;"`
}
