package problem

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

type ProblemInfo struct {
	models.MetaFields
	Slug  string          `json:"slug"`
	Title string          `json:"title"`
	Tags  []*AlgorithmTag `json:"tags"`
}

var ProblemInfoSelection = append([]string{"slug", "title"}, models.MetaFieldsSelection...)

func (p Problem) ToProblemInfo() ProblemInfo {
	return ProblemInfo{
		MetaFields: p.MetaFields,
		Slug:       p.Slug,
		Title:      p.Title,
		Tags:       p.Tags,
	}
}
