package log_module

import (
	"os"
	"runtime"
	"strconv"

	config_module "github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/sirupsen/logrus"
)

const logLevelProp = "log.level"
const logForceQuote = "log.force_quote"
const logTimeOn = "log.time_on"
const logTimeFormat = "log.time_format"

func AppLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"caller": func() string {
			pc := make([]uintptr, 1)
			runtime.Callers(3, pc)
			f := runtime.FuncForPC(pc[0])
			name, line := f.FileLine(pc[0])
			return name + ":" + strconv.Itoa(line)
		}(),
	})
}

func setupLog() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
	lvl := config_module.AppConfig().GetString(logLevelProp)
	forceQuote := config_module.AppConfig().GetBool(logForceQuote)
	fullTimeOn := config_module.AppConfig().GetBool(logTimeOn)
	timestampFormat := config_module.AppConfig().GetString(logTimeFormat)

	logLevel, err := logrus.ParseLevel(lvl)
	if err == nil {
		println("log level:", lvl)
		logrus.SetLevel(logLevel)
	}
	// TODO: control log format in config
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      forceQuote, // value Quote
		FullTimestamp:   fullTimeOn,
		TimestampFormat: timestampFormat,
	})
}

func init() {
	setupLog()
}
