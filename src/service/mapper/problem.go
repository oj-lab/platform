package mapper

import (
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateProblem(tx *gorm.DB, problem model.Problem) error {
	return tx.Create(&problem).Error
}

func GetProblem(tx *gorm.DB, slug string) (*model.Problem, error) {
	db_problem := model.Problem{}
	err := tx.Model(&model.Problem{}).Preload("Tags").Where("Slug = ?", slug).First(&db_problem).Error
	if err != nil {
		return nil, err
	}

	return &db_problem, nil
}

func DeleteProblem(tx *gorm.DB, slug string) error {
	var problem model.Problem
	if err := tx.Where("slug = ?", slug).First(&problem).Error; err != nil {
		return err
	}

	if err := tx.Where("problem_slug = ?", slug).Delete(&model.JudgeTaskSubmission{}).Error; err != nil {
		return err
	}

	return tx.Select(clause.Associations).Delete(&problem).Error
}

func UpdateProblem(tx *gorm.DB, problem model.Problem) error {
	return tx.Model(&model.Problem{Slug: problem.Slug}).Updates(problem).Error
}

type GetProblemOptions struct {
	Selection []string
	Slug      *string
	Title     *string
	Tags      []*model.AlgorithmTag
	Offset    *int
	Limit     *int
}

func buildGetProblemTXByOptions(tx *gorm.DB, options GetProblemOptions, isCount bool) *gorm.DB {
	tagsList := []string{}
	for _, tag := range options.Tags {
		tagsList = append(tagsList, tag.Name)
	}
	tx = tx.Model(&model.Problem{})
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if len(tagsList) > 0 {
		tx = tx.Joins("JOIN problem_algorithm_tags ON problem_algorithm_tags.problem_slug = problems.slug").
			Where("problem_algorithm_tags.algorithm_tag_name in ?", tagsList)
	}
	if options.Slug != nil {
		tx = tx.Where("slug = ?", *options.Slug)
	}
	if options.Title != nil {
		tx = tx.Where("title = ?", *options.Title)
	}
	tx = tx.Distinct().
		Preload("Tags")
	if !isCount {
		if options.Offset != nil {
			tx = tx.Offset(*options.Offset)
		}
		if options.Limit != nil {
			tx = tx.Limit(*options.Limit)
		}
	}

	return tx
}

func CountProblemByOptions(tx *gorm.DB, options GetProblemOptions) (int64, error) {
	var count int64

	tx = buildGetProblemTXByOptions(tx, options, true)
	err := tx.Count(&count).Error

	return count, err
}

func GetProblemListByOptions(tx *gorm.DB, options GetProblemOptions) ([]model.Problem, int64, error) {
	total, err := CountProblemByOptions(tx, options)
	if err != nil {
		return nil, 0, err
	}

	problemList := []model.Problem{}

	tx = buildGetProblemTXByOptions(tx, options, false)
	err = tx.Find(&problemList).Error
	if err != nil {
		return nil, 0, err
	}

	return problemList, total, nil
}

func GetTagsList(problem model.Problem) []string {
	tagsList := []string{}
	for _, tag := range problem.Tags {
		tagsList = append(tagsList, tag.Name)
	}
	return tagsList
}
