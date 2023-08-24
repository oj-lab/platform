package application

import (
	"os"

	"github.com/sirupsen/logrus"
)

const LOG_LEVEL_PROP = "log.level"

func setupLog() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
	lvl := AppConfig.GetString(LOG_LEVEL_PROP)
	logLevel, err := logrus.ParseLevel(lvl)
	if err == nil {
		println("log level:", lvl)
		logrus.SetLevel(logLevel)
	}
}
