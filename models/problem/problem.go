package problem_model

import "github.com/oj-lab/oj-lab-platform/models"

type Problem struct {
	models.MetaFields
	Slug        string          `gorm:"primaryKey" json:"slug"`
	Title       string          `gorm:"not null" json:"title"`
	Description *string         `json:"description,omitempty"`
	Tags        []*AlgorithmTag `gorm:"many2many:problem_algorithm_tags;" json:"tags"`
}

type AlgorithmTag struct {
	models.MetaFields
	Name     string     `gorm:"primaryKey" json:"name"`
	Problems []*Problem `gorm:"many2many:problem_algorithm_tags;" json:"problems,omitempty"`
}

var ProblemInfoSelection = append([]string{"slug", "title"}, models.MetaFieldsSelection...)
