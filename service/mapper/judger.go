package mapper

import (
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/google/uuid"
)

func CreateJudger(judger model.Judger) error {
	judger.UID = uuid.New()
	db := gormAgent.GetDefaultDB()
	return db.Create(&judger).Error
}

func GetJudgerList() ([]model.Judger, error) {
	db := gormAgent.GetDefaultDB()
	var judgers []model.Judger
	err := db.Find(&judgers).Error
	return judgers, err
}
