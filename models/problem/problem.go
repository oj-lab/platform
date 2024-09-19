package problem_model

import "github.com/oj-lab/platform/models"

type Problem struct {
	models.MetaFields
	Slug        string            `json:"slug" gorm:"primaryKey"`
	Title       string            `json:"title" gorm:"not null"`
	Description *string           `json:"description,omitempty"`
	Difficulty  ProblemDifficulty `json:"difficulty,omitempty"`
	Tags        []*ProblemTag     `json:"tags" gorm:"many2many:problem_problem_tags;"`
	Solved      *bool             `json:"solved,omitempty" gorm:"-"`
}

type ProblemDifficulty string

const (
	ProblemDifficultyEasy   ProblemDifficulty = "easy"
	ProblemDifficultyMedium ProblemDifficulty = "medium"
	ProblemDifficultyHard   ProblemDifficulty = "hard"
)

func (d ProblemDifficulty) IsValid() bool {
	switch d {
	case ProblemDifficultyEasy, ProblemDifficultyMedium, ProblemDifficultyHard:
		return true
	}
	return false
}

type ProblemTag struct {
	models.MetaFields
	Name     string     `json:"name" gorm:"primaryKey"`
	Problems []*Problem `json:"problems,omitempty" gorm:"many2many:problem_problem_tags;"`
}
