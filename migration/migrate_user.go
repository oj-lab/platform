package main

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "config/test.ini"
	}
	dataBaseSettings, err := utils.GetDatabaseSettings(configPath)
	if err != nil {
		panic("failed to get database settings")
	}
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
