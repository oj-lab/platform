package model

type DbProblem struct {
	MetaFields
	Slug        string         `gorm:"primaryKey"`
	Title       string         `gorm:"not null"`
	Discription string         `gorm:"not null"`
	ProblemTags []DbProblemTag `gorm:"many2many:problem_tag;"`
}

func (ut DbProblem) TableName() string {
	return "problem"
}

type DbProblemTag struct {
	MetaFields
	Slug     string      `gorm:"primaryKey"`
	Name     string      `gorm:"not null"`
	Problems []DbProblem `gorm:"many2many:problem;"`
}

func (ut DbProblemTag) TableName() string {
	return "problem_tag"
}
