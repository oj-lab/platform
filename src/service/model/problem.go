package model

type Problem struct {
	MetaFields
	Slug        string          `gorm:"primaryKey" json:"slug"`
	Title       string          `gorm:"not null" json:"title"`
	Description *string         `gorm:"not null" json:"description,omitempty"`
	Tags        []*AlgorithmTag `gorm:"many2many:problem_algorithm_tags;" json:"tags"`
}

type AlgorithmTag struct {
	MetaFields
	Name     string     `gorm:"primaryKey" json:"name"`
	Problems []*Problem `gorm:"many2many:problem_algorithm_tags;" json:"problems,omitempty"`
}

type ProblemInfo struct {
	MetaFields
	Slug  string          `json:"slug"`
	Title string          `json:"title"`
	Tags  []*AlgorithmTag `json:"tags"`
}

var ProblemInfoSelection = append([]string{"slug", "title"}, MetaFieldsSelection...)

func (p Problem) ToProblemInfo() ProblemInfo {
	return ProblemInfo{
		MetaFields: p.MetaFields,
		Slug:       p.Slug,
		Title:      p.Title,
		Tags:       p.Tags,
	}
}
