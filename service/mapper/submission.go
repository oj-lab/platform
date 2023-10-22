package mapper

import (
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateSubmission(tx *gorm.DB, submission model.JudgeTaskSubmission) (*model.JudgeTaskSubmission, error) {
	submission.UID = uuid.New()
	submission.MetaFields = model.NewMetaFields()

	return &submission, tx.Create(&submission).Error
}

type GetSubmissionOptions struct {
	Selection      []string
	Statuses       []model.SubmissionStatus
	UserAccount    *string
	ProblemSlug    *string
	Offset         *int
	Limit          *int
	OrderByColumns []model.OrderByColumnOption
}

func BuildGetSubmissionTXByOptions(tx *gorm.DB, options GetSubmissionOptions, isCount bool) *gorm.DB {
	tx = tx.Model(&model.JudgeTaskSubmission{}).
		Preload(clause.Associations)
		// See more in: https://gorm.io/docs/preload.html
		// Preload("User.Roles").Preload("Problem.Tags").Preload(clause.Associations)
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if options.UserAccount != nil {
		tx = tx.Where("user_account = ?", *options.UserAccount)
	}
	if options.ProblemSlug != nil {
		tx = tx.Where("problem_slug = ?", *options.ProblemSlug)
	}
	if len(options.Statuses) > 0 {
		tx = tx.Where("status IN ?", options.Statuses)
	}
	if options.Offset != nil {
		tx = tx.Offset(*options.Offset)
	}
	if options.Limit != nil {
		tx = tx.Limit(*options.Limit)
	}
	for _, option := range options.OrderByColumns {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: option.Column},
			Desc:   option.Desc,
		})
	}

	return tx
}

func GetSubmissionListByOptions(tx *gorm.DB, options GetSubmissionOptions) ([]*model.JudgeTaskSubmission, int64, error) {
	tx = BuildGetSubmissionTXByOptions(tx, options, false)
	var submissions []*model.JudgeTaskSubmission
	err := tx.Find(&submissions).Error
	if err != nil {
		return nil, 0, err
	}

	tx = BuildGetSubmissionTXByOptions(tx, options, true)
	var count int64
	err = tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return submissions, count, nil
}
