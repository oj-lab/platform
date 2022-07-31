package model

import (
	"github.com/OJ-lab/oj-lab-services/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenConnection(settings utils.DatabaseSettings) {
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  utils.GetDatabaseDSN(settings),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
