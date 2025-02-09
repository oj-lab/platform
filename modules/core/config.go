package core_module

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const defaultConfigName = "config"
const defaultOverrideConfigName = "override"
const defaultProjectRootName = "platform"

var Config *viper.Viper

var projectRoot string

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

	Config = viper.GetViper()
	return nil
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
	loadProjectRoot()
	if _, err := os.Stat(projectRoot); err != nil {
		panic(fmt.Sprintf("Project root not found: %v", projectRoot))
	}
	println("Project root:", projectRoot)
	if err := loadConfig(); err != nil {
		panic(fmt.Sprintf("Load config with error: %v", err))
	}
	setupLog()
}
