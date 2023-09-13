package mapper

import (
	"github.com/OJ-lab/oj-lab-services/packages/core"
	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func CreateProblem(problem model.Problem) error {
	db := core.GetDefaultDB()
	return db.Create(&problem).Error
}

func GetProblem(slug string) (*model.Problem, error) {
	db := core.GetDefaultDB()
	db_problem := model.Problem{}
	err := db.Model(&model.Problem{}).Preload("Tags").Where("Slug = ?", slug).First(&db_problem).Error
	if err != nil {
		return nil, err
	}

	return &db_problem, nil
}

func DeleteProblem(problem model.Problem) error {
	db := core.GetDefaultDB()
	return db.Delete(&model.Problem{Slug: problem.Slug}).Error
}

func UpdateProblem(problem model.Problem) error {
	db := core.GetDefaultDB()
	return db.Model(&model.Problem{Slug: problem.Slug}).Updates(problem).Error
}

type GetProblemOptions struct {
	Slug   string
	Title  string
	Tags   []*model.AlgorithmTag
	Offset *int
	Limit  *int
}

func CountProblemByOptions(options GetProblemOptions) (int64, error) {
	db := core.GetDefaultDB()
	var count int64

	tagsList := []string{}
	for _, tag := range options.Tags {
		tagsList = append(tagsList, tag.Slug)
	}

	tx := db.
		Model(&model.Problem{}).
		Joins("JOIN problem_algorithm_tags ON problem_algorithm_tags.problem_slug = problems.slug").
		Where("problem_algorithm_tags.algorithm_tag_slug in ?", tagsList).
		Or("Slug = ?", options.Slug).
		Or("Title = ?", options.Title).
		Distinct().
		Preload("Tags")

	err := tx.Count(&count).Error

	return count, err
}

func GetProblemByOptions(options GetProblemOptions) ([]model.Problem, int64, error) {
	total, err := CountProblemByOptions(options)
	if err != nil {
		return nil, 0, err
	}

	db := core.GetDefaultDB()
	db_problems := []model.Problem{}
	tagsList := []string{}
	for _, tag := range options.Tags {
		tagsList = append(tagsList, tag.Slug)
	}
	tx := db.
		Model(&model.Problem{}).
		Joins("JOIN problem_algorithm_tags ON problem_algorithm_tags.problem_slug = problems.slug").
		Where("problem_algorithm_tags.algorithm_tag_slug in ?", tagsList).
		Or("Slug = ?", options.Slug).
		Or("Title = ?", options.Title).
		Distinct().
		Preload("Tags")

	if options.Offset != nil {
		tx = tx.Offset(*options.Offset)
	}
	if options.Limit != nil {
		tx = tx.Limit(*options.Limit)
	}

	err = tx.Find(&db_problems).Error
	if err != nil {
		return nil, 0, err
	}

	return db_problems, total, nil
}

func GetTagsList(problem model.Problem) []string {
	tagsList := []string{}
	for _, tag := range problem.Tags {
		tagsList = append(tagsList, tag.Slug)
	}
	return tagsList
}
