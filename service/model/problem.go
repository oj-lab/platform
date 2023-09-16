package model

type Problem struct {
	MetaFields
	Slug        string          `gorm:"primaryKey"`
	Title       string          `gorm:"not null"`
	Description string          `gorm:"not null"`
	Tags        []*AlgorithmTag `gorm:"many2many:problem_algorithm_tags;"`
}

type AlgorithmTag struct {
	MetaFields
	Slug     string     `gorm:"primaryKey"`
	Name     string     `gorm:"not null"`
	Problems []*Problem `gorm:"many2many:problem_algorithm_tags;"`
}
