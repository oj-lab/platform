package model

type Problem struct {
	MetaFields
	Slug        string          `gorm:"primaryKey" json:"slug"`
	Title       string          `gorm:"not null" json:"title"`
	Description *string         `gorm:"not null" json:"description,omitempty"`
	Tags        []*AlgorithmTag `gorm:"many2many:problem_algorithm_tags;" json:"tags"`
}

var ProblemInfoSelection = append([]string{"slug", "title"}, MetaFieldsSelection...)

type AlgorithmTag struct {
	MetaFields
	Slug     string     `gorm:"primaryKey" json:"slug"`
	Name     string     `gorm:"not null" json:"name"`
	Problems []*Problem `gorm:"many2many:problem_algorithm_tags;" json:"problems,omitempty"`
}
