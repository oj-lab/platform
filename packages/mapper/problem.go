package mapper

import (
	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func CreateProblem(problem model.Problem) error {
	db := application.GetDefaultDB()
	return db.Create(&problem).Error
}

func GetProblem(slug string) (model.Problem, error) {
	db := application.GetDefaultDB()
	db_problem := model.Problem{}
	err := db.Model(&model.Problem{}).Preload("Tags").Where("Slug = ?", slug).First(&db_problem).Error
	return db_problem, err
}

func DeleteProblem(problem model.Problem) error {
	db := application.GetDefaultDB()
	return db.Delete(&model.Problem{Slug: problem.Slug}).Error
}

func UpdateProblem(problem model.Problem) error {
	db := application.GetDefaultDB()
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
	db := application.GetDefaultDB()
	var count int64

	tx := db.
		Model(&model.Problem{}).
		Where("Slug = ?", options.Slug).
		Or("Title = ?", options.Title).
		Or("Tags in ?", options.Tags)

	err := tx.Count(&count).Error

	return count, err
}

func GetProblemByOptions(options GetProblemOptions) ([]model.Problem, int64, error) {
	total, err := CountProblemByOptions(options)
	if err != nil {
		return nil, 0, err
	}

	db := application.GetDefaultDB()
	db_problems := []model.Problem{}
	tx := db.
		Where("Slug = ?", options.Slug).
		Or("Title = ?", options.Title).
		Or("Tags in ?", options.Tags)
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
