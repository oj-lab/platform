package main

import (
	"github.com/OJ-lab/oj-lab-services/user-service/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	dataBaseSettings := utils.GetDatabaseSettings(os.Args[1])
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  utils.GetDatabaseDSN(dataBaseSettings),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to migrate database")
	}
}
