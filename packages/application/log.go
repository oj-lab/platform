package application

import (
	"os"

	"github.com/sirupsen/logrus"
)

const logLevelProp = "log.level"

func setupLog() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
	lvl := AppConfig.GetString(logLevelProp)
	logLevel, err := logrus.ParseLevel(lvl)
	if err == nil {
		println("log level:", lvl)
		logrus.SetLevel(logLevel)
	}
}
