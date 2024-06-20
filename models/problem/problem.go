package problem

import "github.com/oj-lab/oj-lab-platform/models"

type Problem struct {
	models.MetaFields
	Slug        string      `json:"slug" gorm:"primaryKey"`
	Info        ProblemInfo `json:"info"`
	Description *string     `json:"description,omitempty"`
}

var InfoSelection = append([]string{"slug", "title"}, models.MetaFieldsSelection...)
