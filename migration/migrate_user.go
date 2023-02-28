package main

import (
	"os"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "config/ini/test.ini"
	}
	dataBaseSettings, err := config.GetDatabaseSettings(configPath)
	if err != nil {
		panic("failed to get database settings")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  utils.MustGetDatabaseDSN(*dataBaseSettings),
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
