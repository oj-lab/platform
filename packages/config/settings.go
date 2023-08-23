package config

import (
	"log"

	"gopkg.in/ini.v1"
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
	AuthOn    bool
	Port      string
	Mode      string
	CookieAge string
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

func GetDatabaseSettings(source interface{}) (*DatabaseSettings, error) {
	cfg, err := ini.Load(source)
	if err != nil {
		return nil, err
	}

	var databaseSettings DatabaseSettings
	err = cfg.Section("database").MapTo(&databaseSettings)
	if err != nil {
		return nil, err
	}
	log.Println("load databaseSettings")
	return &databaseSettings, nil
}
