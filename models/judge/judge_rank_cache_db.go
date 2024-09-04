package judge_model

import (
	"errors"

	"github.com/oj-lab/oj-lab-platform/models"
	"gorm.io/gorm"
)

func CreateJudgeRankCache(tx *gorm.DB, rankCache JudgeRankCache) (*JudgeRankCache, error) {
	rankCache.MetaFields = models.NewMetaFields()
	return &rankCache, tx.Create(&rankCache).Error
}

func DeleteJudgeRankCache(tx *gorm.DB, userAccount string) error {
	var judgeRankCache JudgeRankCache
	if err := tx.Where("user_account = ?", userAccount).First(&judgeRankCache).Error; err != nil {
		return err
	}
	return tx.Delete(&judgeRankCache).Error
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

func (rankCache *JudgeRankCache) BeforeSave(tx *gorm.DB) (err error) {
	if rankCache.Points > rankCache.TotalSubmissions {
		return errors.New("TotalSubmissions must >= Points")
	}
	return nil
}
