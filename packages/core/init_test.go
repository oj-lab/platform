package core

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInit(T *testing.T) {
	logLevel := viper.GetString(logLevelProp)
	if logLevel != "debug" {
		T.Errorf("log level is not debug but %s", logLevel)
	}
}
