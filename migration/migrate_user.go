package main

import (
	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/sirupsen/logrus"
)

func main() {
	db := application.GetDefaultDB()
	err := db.AutoMigrate(&model.UserTable{})
	if err != nil {
		panic("failed to migrate database")
	}
	logrus.Info("migrate user table success")
}
