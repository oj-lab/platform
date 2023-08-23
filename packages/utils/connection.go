package utils

import (
	"fmt"
	"log"

	"github.com/OJ-lab/oj-lab-services/packages/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustGetDBConnection(settings config.DatabaseSettings) *gorm.DB {
	var err error
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  MustGetDatabaseDSN(settings),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func MustCreateDatabase(settings config.DatabaseSettings) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  MustGetPSqlDSN(settings),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect psql")
	}
	var exists bool
	checkDatabaseCmd := fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_database WHERE datname = '%s') AS FOUND;", settings.DbName)
	rs := db.Raw(checkDatabaseCmd).Scan(&exists)
	if rs.Error != nil {
		panic("failed to check database")
	}
	if !exists {
		createDatabaseCmd := fmt.Sprintf("CREATE DATABASE %s;", settings.DbName)
		rs := db.Exec(createDatabaseCmd)
		if rs.Error != nil {
			panic("failed to create database")
		} else {
			log.Printf("database %s created\n", settings.DbName)
		}
	} else {
		log.Printf("database %s existed\n", settings.DbName)
	}
}

func MustGetPSqlDSN(settings config.DatabaseSettings) string {
	return "user=" + settings.User + " password=" + settings.Password + " host=" + settings.Host + " port=" + settings.Port + " sslmode=disable TimeZone=Asia/Shanghai"
}

func MustGetDatabaseDSN(settings config.DatabaseSettings) string {
	return "user=" + settings.User + " password=" + settings.Password + " dbname=" + settings.DbName + " host=" + settings.Host + " port=" + settings.Port + " sslmode=disable TimeZone=Asia/Shanghai"
}
