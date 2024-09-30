package judge_model

import (
	"github.com/google/uuid"
	"github.com/oj-lab/platform/models"
	"gorm.io/gorm"
)

func CreateJudgeResult(tx *gorm.DB, result JudgeResult) (*JudgeResult, error) {
	result.UID = uuid.New()
	result.MetaFields = models.NewMetaFields()
	if !result.Verdict.IsValid() {
		return nil, ErrInvalidJudgeVerdict
	}

	return &result, tx.Create(&result).Error
}

func DeleteJudgeResultByJudgeUID(tx *gorm.DB, judgeUID uuid.UUID) error {
	return tx.Where("judge_uid = ?", judgeUID).Delete(&JudgeResult{}).Error
}
