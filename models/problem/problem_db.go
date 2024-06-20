package problem

import (
	"gorm.io/gorm"
)

func CreateProblem(tx *gorm.DB, problem Problem) error {
	return tx.Create(&problem).Error
}

func GetProblem(tx *gorm.DB, slug string) (*Problem, error) {
	db_problem := Problem{}
	err := tx.Model(&Problem{}).Preload("Info").
		Preload("Info.Tags").
		Where("slug = ?", slug).First(&db_problem).Error
	if err != nil {
		return nil, err
	}

	return &db_problem, nil
}

func DeleteProblem(tx *gorm.DB, slug string) error {
	err := tx.Model(&ProblemInfo{
		ProblemSlug: slug,
	}).Association("Tags").Clear()
	if err != nil {
		return err
	}

	return tx.Select("Info").
		Delete(&Problem{Slug: slug}).Error
}

func UpdateProblem(tx *gorm.DB, problem Problem) error {
	return tx.Model(&Problem{Slug: problem.Slug}).Updates(problem).Error
}

type GetProblemOptions struct {
	Slug   *string
	Title  *string
	Tags   []*ProblemTag
	Offset *int
	Limit  *int
}

func GetTagsList(problem Problem) []string {
	tagsList := []string{}
	for _, tag := range problem.Info.Tags {
		tagsList = append(tagsList, tag.Name)
	}
	return tagsList
}
