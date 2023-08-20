package config

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

func (se ServiceEnv) isValid() bool {
	if se == DEV_SERVICE_ENV || se == PRD_SERVICE_ENV {
		return true
	}
	return false
}

func loadConfig(basePath string) error {
	viper.AddConfigPath(basePath)

	env := os.Getenv(SERVICE_ENV_KEY)
	if ServiceEnv(env).isValid() {
		viper.SetConfigName(env)
	} else {
		viper.SetConfigName(string(DEV_SERVICE_ENV))
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetConfigName(OVERRIDE_CONFIG_NAME)
	err = viper.MergeInConfig()
	if err == nil {
		println("Found override config, merging...")
	}

	return nil
}
