package judge_model

import (
	"github.com/oj-lab/oj-lab-platform/models"
	"gorm.io/gorm"
)

func CreateJudgeScoreCache(tx *gorm.DB, scoreCache JudgeScoreCache) (*JudgeScoreCache, error) {
	scoreCache.MetaFields = models.NewMetaFields()
	return &scoreCache, tx.Create(&scoreCache).Error
}

func GetJudgeScoreCache(tx *gorm.DB, userAccount string, problemSlug string) (*JudgeScoreCache, error) {
	scoreCache := JudgeScoreCache{}
	err := tx.Model(&JudgeScoreCache{}).
		Where("user_account = ?", userAccount).
		Where("problem_slug = ?", problemSlug).
		First(&scoreCache).Error
	if err != nil {
		return nil, err
	}
	return &scoreCache, nil
}

func UpdateJudgeScoreCache(tx *gorm.DB, scoreCache JudgeScoreCache) error {
	return tx.Model(&scoreCache).Updates(scoreCache).Error
}
