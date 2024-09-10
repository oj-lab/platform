package judge_model

import (
	user_model "github.com/oj-lab/platform/models/user"
	"gorm.io/gorm"
)

type GetRankCacheOptions struct {
	Selection   []string
	UserAccount *string
	Offset      *int
	Limit       *int
}

func buildGetRankCacheTXByOptions(tx *gorm.DB, options GetRankCacheOptions, isCount bool) *gorm.DB {
	tx = tx.Model(&JudgeRankCache{}).Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select(user_model.PublicUserSelection)
	})
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if options.UserAccount != nil {
		tx = tx.Where("user_account = ?", *options.UserAccount)
	}

	if !isCount {
		if options.Offset != nil {
			tx = tx.Offset(*options.Offset)
		}
		if options.Limit != nil {
			tx = tx.Limit(*options.Limit)
		}
	}
	tx = tx.Order("points DESC, total_submissions ASC")
	return tx
}

func CountRankByOptions(tx *gorm.DB, options GetRankCacheOptions) (int64, error) {
	var count int64

	tx = buildGetRankCacheTXByOptions(tx, options, true)
	err := tx.Count(&count).Error

	return count, err
}

func GetRankCacheListByOptions(tx *gorm.DB, options GetRankCacheOptions) ([]JudgeRankCache, error) {
	rankInfoList := []JudgeRankCache{}
	tx = buildGetRankCacheTXByOptions(tx, options, false)
	err := tx.Find(&rankInfoList).Error
	if err != nil {
		return nil, err
	}
	return rankInfoList, nil
}
