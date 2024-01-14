package mapper

import (
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"gorm.io/gorm"
)

func CreateJudger(tx *gorm.DB, judger model.Judger) error {
	return tx.Create(&judger).Error
}

func GetJudgerList(tx *gorm.DB) ([]model.Judger, error) {
	var judgers []model.Judger
	err := tx.Find(&judgers).Error
	return judgers, err
}
