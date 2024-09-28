package log_module

import (
	"os"
	"runtime"
	"strconv"

	config_module "github.com/oj-lab/platform/modules/config"
	"github.com/sirupsen/logrus"
)

const logLevelProp = "log.level"
const logFormatProp = "log.format"
const logPrettyJson = "log.pretty_json"
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
	format := config_module.AppConfig().GetString(logFormatProp)
	lvl := config_module.AppConfig().GetString(logLevelProp)
	prettyJson := config_module.AppConfig().GetBool(logPrettyJson)
	timeOn := config_module.AppConfig().GetBool(logTimeOn)
	timestampFormat := config_module.AppConfig().GetString(logTimeFormat)

	logLevel, err := logrus.ParseLevel(lvl)
	if err == nil {
		println("log level:", lvl)
		logrus.SetLevel(logLevel)
	}
	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  timestampFormat,
			DisableTimestamp: !timeOn,
			PrettyPrint:      prettyJson,
		})
	}
}

func init() {
	setupLog()
}
