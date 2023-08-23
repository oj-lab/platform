package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const LOG_LEVEL_PROP = "log.level"

func SetupLog() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
	logLevel, err := logrus.ParseLevel(viper.GetString(LOG_LEVEL_PROP))
	if err == nil {
		logrus.SetLevel(logLevel)
	}
}
