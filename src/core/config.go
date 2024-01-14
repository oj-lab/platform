package core

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

const SERVICE_ENV_KEY = "OJ_LAB_SERVICE_ENV"
const PROJECT_ROOT_ENV_KEY = "OJ_LAB_PROJECT_ROOT"
const OVERRIDE_CONFIG_NAME_ENV_KEY = "OJ_LAB_OVERRIDE_CONFIG_NAME"

const DEFAULT_OVERRIDE_CONFIG_NAME = "override"
const DEFAULT_PROJECT_ROOT = "oj-lab-platform"

type ServiceEnv string

const (
	DEV_SERVICE_ENV ServiceEnv = "development"
	PRD_SERVICE_ENV ServiceEnv = "production"
)

var serviceEnv ServiceEnv
var AppConfig *viper.Viper

func (se ServiceEnv) isValid() bool {
	if se == DEV_SERVICE_ENV || se == PRD_SERVICE_ENV {
		return true
	}
	return false
}

func IsDevEnv() bool {
	return serviceEnv == DEV_SERVICE_ENV
}

func loadConfig(basePath string) error {
	viper.AddConfigPath(basePath)

	serviceEnv = DEV_SERVICE_ENV
	env := os.Getenv(SERVICE_ENV_KEY)
	if ServiceEnv(env).isValid() {
		serviceEnv = ServiceEnv(env)
	}
	println("Env:", serviceEnv)
	viper.SetConfigName(string(serviceEnv))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	overrideConfigName := os.Getenv(OVERRIDE_CONFIG_NAME_ENV_KEY)
	if overrideConfigName == "" {
		overrideConfigName = DEFAULT_OVERRIDE_CONFIG_NAME
	}
	println("Set override config name:", overrideConfigName)

	viper.SetConfigName(overrideConfigName)
	err = viper.MergeInConfig()
	if err == nil {
		println("Found override config, merged")
	}

	AppConfig = viper.GetViper()
	return nil
}

func GetProjectRoot() string {
	projectRoot := os.Getenv(PROJECT_ROOT_ENV_KEY)
	if projectRoot == "" {
		projectRoot = DEFAULT_PROJECT_ROOT
	}

	wd, err := os.Getwd()
	if err != nil {
		panic("Cannot find working dir")
	}
	isRoot := path.Base(wd) == projectRoot
	for !isRoot && wd != "/" {
		wd = path.Dir(wd)
		isRoot = path.Base(wd) == projectRoot
	}
	if wd == "/" {
		panic("Cannot find project root folder")
	}
	return wd
}

func init() {
	projectRoot := GetProjectRoot()
	println("Initing config with project root:", projectRoot)
	loadConfig(path.Join(projectRoot, "environment/configs"))
	setupLog()
}
