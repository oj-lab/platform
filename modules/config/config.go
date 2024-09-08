package config_module

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const serviceEnvEnvKey = "OJ_LAB_SERVICE_ENV"

const defaultConfigName = "config"
const defaultOverrideConfigName = "override"
const defaultProjectRootName = "platform"

type ServiceEnv string

const (
	serviceEnvDev ServiceEnv = "development"
	serviceEnvPrd ServiceEnv = "production"
)

var serviceEnv ServiceEnv
var projectRoot string
var appConfig *viper.Viper

func (se ServiceEnv) isValid() bool {
	if se == serviceEnvDev || se == serviceEnvPrd {
		return true
	}
	return false
}

func IsDevEnv() bool {
	return serviceEnv == serviceEnvDev
}

func loadServiceEnv() {
	serviceEnv = serviceEnvDev
	env := os.Getenv(serviceEnvEnvKey)
	if ServiceEnv(env).isValid() {
		serviceEnv = ServiceEnv(env)
	}
	println("Env:", serviceEnv)
}

func loadConfig() error {
	viper.AddConfigPath(projectRoot)

	viper.SetConfigName(defaultConfigName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetConfigName(defaultOverrideConfigName)
	err = viper.MergeInConfig()
	if err == nil {
		println("Found override config, merged")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	appConfig = viper.GetViper()
	return nil
}

func AppConfig() *viper.Viper {
	return appConfig
}

func loadProjectRoot() {
	// Try to locate project root, then find the workdir
	wd, err := os.Getwd()
	if err != nil {
		panic("Cannot get cwd")
	}
	println("Checking workdir from cwd:", wd)
	isProjectRoot := path.Base(wd) == defaultProjectRootName
	for !isProjectRoot && wd != "/" {
		wd = path.Dir(wd)
		isProjectRoot = path.Base(wd) == defaultProjectRootName
	}
	if wd == "/" {
		panic("Cannot find projectRoot")
	}
	projectRoot = wd
}

func ProjectRoot() string {
	return projectRoot
}

func init() {
	loadServiceEnv()
	loadProjectRoot()
	if _, err := os.Stat(projectRoot); err != nil {
		panic(fmt.Sprintf("Project root not found: %v", projectRoot))
	}
	println("Project root:", projectRoot)
	if err := loadConfig(); err != nil {
		panic(fmt.Sprintf("Load config with error: %v", err))
	}
}
