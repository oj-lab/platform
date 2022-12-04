package utils

import (
	"github.com/OJ-lab/oj-lab-services/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection(settings config.DatabaseSettings) *gorm.DB {
	var err error
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDatabaseDSN(settings),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func GetDatabaseDSN(settings config.DatabaseSettings) string {
	return "user=" + settings.User + " password=" + settings.Password + " dbname=" + settings.DbName + " host=" + settings.Host + " port=" + settings.Port + " sslmode=disable TimeZone=Asia/Shanghai"
}
