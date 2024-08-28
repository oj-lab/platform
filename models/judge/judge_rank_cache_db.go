package judge_model

import (
	"github.com/oj-lab/oj-lab-platform/models"
	"gorm.io/gorm"
)

func CreateJudgeRankCache(tx *gorm.DB, rankCache JudgeRankCache) (*JudgeRankCache, error) {
	rankCache.MetaFields = models.NewMetaFields()
	return &rankCache, tx.Create(&rankCache).Error
}

func GetJudgeRankCache(tx *gorm.DB, userAccount string) (*JudgeRankCache, error) {
	rankCache := JudgeRankCache{}
	err := tx.Model(&JudgeRankCache{}).
		Where("user_account = ?", userAccount).
		First(&rankCache).Error
	if err != nil {
		return nil, err
	}
	return &rankCache, nil
}

func UpdateJudgeRankCache(tx *gorm.DB, rankCache JudgeRankCache) error {
	return tx.Model(&rankCache).Updates(rankCache).Error
}
