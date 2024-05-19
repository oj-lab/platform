package log

import (
	"os"
	"runtime"

	"github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/sirupsen/logrus"
)

const logLevelProp = "log.level"

func AppLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"CALLER": func() string {
			pc := make([]uintptr, 1)
			runtime.Callers(3, pc)
			f := runtime.FuncForPC(pc[0])
			return f.Name()
		}(),
	})
}

func setupLog() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
	lvl := config.AppConfig.GetString(logLevelProp)
	logLevel, err := logrus.ParseLevel(lvl)
	if err == nil {
		println("log level:", lvl)
		logrus.SetLevel(logLevel)
	}
}

func init() {
	setupLog()
}
