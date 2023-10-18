package mapper

import (
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func CreateJudger(judger model.Judger) error {
	db := gormAgent.GetDefaultDB()
	return db.Create(&judger).Error
}

func GetJudgerList() ([]model.Judger, error) {
	db := gormAgent.GetDefaultDB()
	var judgers []model.Judger
	err := db.Find(&judgers).Error
	return judgers, err
}
