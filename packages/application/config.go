package app

import (
	"os"

	"github.com/spf13/viper"
)

const SERVICE_ENV_KEY = "OJ_LAB_SERVICE_ENV"
const OVERRIDE_CONFIG_NAME = "override"

type ServiceEnv string

const (
	DEV_SERVICE_ENV ServiceEnv = "development"
	PRD_SERVICE_ENV ServiceEnv = "production"
)

var serviceEnv ServiceEnv

func (se ServiceEnv) isValid() bool {
	if se == DEV_SERVICE_ENV || se == PRD_SERVICE_ENV {
		return true
	}
	return false
}

func IsDevEnv() bool {
	return serviceEnv == DEV_SERVICE_ENV
}

func LoadConfig(basePath string) error {
	viper.AddConfigPath(basePath)

	serviceEnv = DEV_SERVICE_ENV
	env := os.Getenv(SERVICE_ENV_KEY)
	if ServiceEnv(env).isValid() {
		serviceEnv = ServiceEnv(env)
	}
	viper.SetConfigName(string(serviceEnv))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetConfigName(OVERRIDE_CONFIG_NAME)
	err = viper.MergeInConfig()
	if err == nil {
		println("Found override config, merged")
	}

	return nil
}
