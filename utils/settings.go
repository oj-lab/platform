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

type JWTSettings struct {
	Secret   string
	Duration string
}

type ServiceSettings struct {
	AuthOn bool
	Port   string
	Mode   string
}

func GetServiceSettings(source interface{}) (ServiceSettings, error) {
	var cfg *ini.File
	cfg, _ = ini.Load(source)
	var serviceSettings ServiceSettings
	err := cfg.Section("service").MapTo(&serviceSettings)
	if err != nil {
		return ServiceSettings{}, err
	}
	log.Println("load serviceSettings")
	return serviceSettings, nil
}

func GetJWTSettings(source interface{}) (JWTSettings, error) {
	var cfg *ini.File
	cfg, _ = ini.Load(source)
	var jwtSettings JWTSettings
	err := cfg.Section("jwt").MapTo(&jwtSettings)
	if err != nil {
		return JWTSettings{}, err
	}
	log.Println("load jwtSettings")
	return jwtSettings, nil
}

func GetDatabaseSettings(source interface{}) (DatabaseSettings, error) {
	var cfg *ini.File
	cfg, _ = ini.Load(source)
	var databaseSettings DatabaseSettings
	err := cfg.Section("database").MapTo(&databaseSettings)
	if err != nil {
		return DatabaseSettings{}, err
	}
	log.Println("load databaseSettings")
	return databaseSettings, nil
}

func GetDatabaseDSN(settings DatabaseSettings) string {
	return "user=" + settings.User + " password=" + settings.Password + " dbname=" + settings.DbName + " host=" + settings.Host + " port=" + settings.Port + " sslmode=disable TimeZone=Asia/Shanghai"
}
