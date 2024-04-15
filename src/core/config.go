package core

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const serviceEnvEnvKey = "OJ_LAB_SERVICE_ENV"
const workdirEnvKey = "OJ_LAB_WORKDIR"

const defaultConfigName = "config"
const defaultOverrideConfigName = "override"
const defaultProjectRootName = "oj-lab-platform"
const defaultProjectWorkdirFolder = "workdirs"

type ServiceEnv string

const (
	serviceEnvDev ServiceEnv = "development"
	serviceEnvPrd ServiceEnv = "production"
)

var serviceEnv ServiceEnv
var Workdir string
var AppConfig *viper.Viper

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
	viper.AddConfigPath(Workdir)

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

	AppConfig = viper.GetViper()
	return nil
}

func loadWorkdir() {
	Workdir = os.Getenv(workdirEnvKey)
	if Workdir != "" {
		return
	}

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
	Workdir = path.Join(wd, defaultProjectWorkdirFolder, string(serviceEnv))
}

func init() {
	loadServiceEnv()
	loadWorkdir()
	if _, err := os.Stat(Workdir); err != nil {
		panic(fmt.Sprintf("Set workdir %s with error: %v", Workdir, err))
	}
	println("Workdir:", Workdir)
	if err := loadConfig(); err != nil {
		panic(fmt.Sprintf("Load config with error: %v", err))
	}
	setupLog()
}
