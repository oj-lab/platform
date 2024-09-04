package judge_model

import "gorm.io/gorm"

type GetRankOptions struct {
	Selection   []string
	UserAccount *string
	Offset      *int
	Limit       *int
}

func buildGetRankCacheTXByOptions(tx *gorm.DB, options GetRankOptions, isCount bool) *gorm.DB {
	tx = tx.Model(&JudgeRankCache{})
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if options.UserAccount != nil {
		tx = tx.Where("user_account = ?", *options.UserAccount)
	}

	tx = tx.Distinct().Preload("User")
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

func GetRankCacheListByOptions(tx *gorm.DB, options GetRankOptions) ([]JudgeRankCache, error) {
	rankInfoList := []JudgeRankCache{}
	tx = buildGetRankCacheTXByOptions(tx, options, false)
	err := tx.Find(&rankInfoList).Error
	if err != nil {
		return nil, err
	}
	return rankInfoList, nil
}
