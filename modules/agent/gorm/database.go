package gormAgent

import (
	"github.com/oj-lab/oj-lab-platform/modules/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsnProp = "database.dsn"

var db *gorm.DB

var dsn string

func init() {
	dsn = config.AppConfig.GetString(dsnProp)
	if dsn == "" {
		panic("database dsn is not set")
	}
}

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
