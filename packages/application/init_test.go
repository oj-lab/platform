package application

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInit(T *testing.T) {
	logLevel := viper.GetString(LOG_LEVEL_PROP)
	if logLevel != "debug" {
		T.Errorf("log level is not debug but %s", logLevel)
	}
}
