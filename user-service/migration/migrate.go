package main

import (
	sharedTools "github.com/OJ-lab/oj-lab-services/shared-tools"
	"github.com/OJ-lab/oj-lab-services/user-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	dataBaseSettings := sharedTools.GetDatabaseSettings(os.Args[1])
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  sharedTools.GetDatabaseDSN(dataBaseSettings),
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
