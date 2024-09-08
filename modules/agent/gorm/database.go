package gorm_agent

import (
	config_module "github.com/oj-lab/platform/modules/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsnProp = "database.dsn"

var db *gorm.DB

var dsn string

func init() {
	dsn = config_module.AppConfig().GetString(dsnProp)
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
		}), &gorm.Config{
			Logger: getLogger(),
		})
		if err != nil {
			panic("failed to connect database")
		}
	}

	return db
}
