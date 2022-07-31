package utils

import (
	"gopkg.in/ini.v1"
	"log"
)

type DatabaseSettings struct {
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

func GetDatabaseSettings(source interface{}) DatabaseSettings {
	var cfg *ini.File
	cfg, _ = ini.Load(source)
	var databaseSettings DatabaseSettings
	err := cfg.Section("database").MapTo(&databaseSettings)
	if err != nil {
		return DatabaseSettings{}
	}
	log.Print(databaseSettings)
	return databaseSettings
}

func GetDatabaseDSN(settings DatabaseSettings) string {
	return "user=" + settings.User + " password=" + settings.Password + " dbname=" + settings.DbName + " host=" + settings.Host + " port=" + settings.Port + " sslmode=disable TimeZone=Asia/Shanghai"
}
