package mapper

import (
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateSubmission(submission model.JudgeTaskSubmission) (*model.JudgeTaskSubmission, error) {
	submission.UID = uuid.New()
	submission.MetaFields = model.NewMetaFields()
	db := gormAgent.GetDefaultDB()

	return &submission, db.Create(&submission).Error
}

type GetSubmissionOptions struct {
	Selection      []string
	UserAccount    *string
	ProblemSlug    *string
	Offset         *int
	Limit          *int
	OrderByColumns []model.OrderByColumnOption
}

func buildGetSubmissionTXByOptions(db *gorm.DB, options GetSubmissionOptions, isCount bool) *gorm.DB {
	tx := db.Model(&model.JudgeTaskSubmission{}).Preload("User").Preload("Problem")
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if options.UserAccount != nil {
		tx = tx.Where("user_account = ?", *options.UserAccount)
	}
	if options.ProblemSlug != nil {
		tx = tx.Where("problem_slug = ?", *options.ProblemSlug)
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

func GetSubmissionListByOptions(options GetSubmissionOptions) ([]*model.JudgeTaskSubmission, int64, error) {
	db := gormAgent.GetDefaultDB()
	tx := buildGetSubmissionTXByOptions(db, options, false)
	var submissions []*model.JudgeTaskSubmission
	err := tx.Find(&submissions).Error
	if err != nil {
		return nil, 0, err
	}

	tx = buildGetSubmissionTXByOptions(db, options, true)
	var count int64
	err = tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return submissions, count, nil
}
