package shared_tools

import (
	"gopkg.in/ini.v1"
	"log"
)

type DatabaseSettings struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
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
