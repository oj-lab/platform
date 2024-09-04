package judge_model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/models"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateJudge(tx *gorm.DB, judge Judge) (*Judge, error) {
	judge.UID = uuid.New()
	judge.MetaFields = models.NewMetaFields()
	if judge.UserAccount == "" {
		judge.UserAccount = "anonymous"
	}

	return &judge, tx.Create(&judge).Error
}

// only count if when accept && accept time < SolvedTime
func GetBeforeSubmission(tx *gorm.DB, judge Judge) (int, error) {
	var count int64
	err := tx.Model(&Judge{}).
		Where("create_at < ?", judge.CreateAt).
		Where("status = ?", JudgeStatusFinished).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return int(count + 1), err
}

func GetJudge(tx *gorm.DB, uid uuid.UUID) (*Judge, error) {
	judge := Judge{}
	err := tx.Model(&Judge{}).
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select(user_model.PublicUserSelection)
		}).
		Preload("Problem", func(tx *gorm.DB) *gorm.DB {
			return tx.Select(problem_model.ProblemInfoSelection)
		}).
		Preload("Results").
		Where("UID = ?", uid).First(&judge).Error
	if err != nil {
		return nil, err
	}
	return &judge, nil
}

type GetJudgeOptions struct {
	Selection      []string
	UserAccount    *string
	ProblemSlugs   []string
	Statuses       []JudgeStatus
	Verdicts       []JudgeVerdict
	Offset         *int
	Limit          *int
	OrderByColumns []models.OrderByColumnOption
}

func buildGetJudgesTXByOptions(
	tx *gorm.DB, options GetJudgeOptions, isCount bool,
) *gorm.DB {
	tx = tx.Model(&Judge{}).
		Preload(clause.Associations) // See more in: https://gorm.io/docs/preload.html
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if options.UserAccount != nil {
		tx = tx.Where("user_account = ?", *options.UserAccount)
	}
	if len(options.Statuses) > 0 {
		tx = tx.Where("status IN ?", options.Statuses)
	}
	if len(options.Verdicts) > 0 {
		tx = tx.Where("verdict IN ?", options.Verdicts)
	}
	if len(options.ProblemSlugs) > 0 {
		tx = tx.Where("problem_slug IN ?", options.ProblemSlugs)
	}

	if !isCount {
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
	}

	return tx
}

func CountJudgeByOptions(tx *gorm.DB, options GetJudgeOptions) (int64, error) {
	tx = buildGetJudgesTXByOptions(tx, options, true)
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetJudgeListByOptions(
	tx *gorm.DB, options GetJudgeOptions,
) ([]*Judge, error) {
	tx = buildGetJudgesTXByOptions(tx, options, false)
	var judges []*Judge
	err := tx.Find(&judges).Error
	if err != nil {
		return nil, err
	}

	return judges, nil
}

func UpdateJudge(tx *gorm.DB, judge Judge) error {
	updatingJudge := Judge{}
	if judge.UID != uuid.Nil {
		err := tx.Where("uid = ?", judge.UID).First(&updatingJudge).Error
		if err != nil {
			return err
		}
	} else if judge.RedisStreamID != "" {
		err := tx.Where("redis_stream_id = ?", judge.RedisStreamID).
			First(&updatingJudge).Error
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("judge uid and redis stream id are both empty")
	}

	if judge.Status != "" {
		updatingJudge.Status = judge.Status
	}
	if judge.RedisStreamID != "" {
		updatingJudge.RedisStreamID = judge.RedisStreamID
	}
	if judge.ResultCount != 0 {
		updatingJudge.ResultCount = judge.ResultCount
	}
	if judge.Verdict != "" {
		updatingJudge.Verdict = judge.Verdict
	}
	if judge.MetaFields.CreateAt != nil {
		updatingJudge.MetaFields = judge.MetaFields
	}

	return tx.Model(&updatingJudge).Updates(updatingJudge).Error
}
