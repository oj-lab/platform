package application

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var dsn string

func GetDefaultDB() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}

	return db
}

func init() {
	dsn = AppConfig.GetString("database.dsn")
	if dsn == "" {
		panic("database dsn is not set")
	}
}
