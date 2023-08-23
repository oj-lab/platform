package main

import (
	"os"

	"github.com/OJ-lab/oj-lab-services/packages/config"
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/OJ-lab/oj-lab-services/packages/utils"
)

func main() {
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "packages/config/ini/example.ini"
	}
	dataBaseSettings, err := config.GetDatabaseSettings(configPath)
	if err != nil {
		panic("failed to get database settings")
	}
	utils.MustCreateDatabase(*dataBaseSettings)
	db := utils.MustGetDBConnection(*dataBaseSettings)
	err = db.AutoMigrate(&model.UserTable{})
	if err != nil {
		panic("failed to migrate database")
	}
}
