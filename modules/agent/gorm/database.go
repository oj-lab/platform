package gorm_agent

import (
	"log/slog"

	core_module "github.com/oj-lab/platform/modules/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dsnConfigKey                = "gorm.dsn"
	loggerIgnoreRunLogConfigKey = "gorm.logger.ignore_run_log"
)

var (
	db  *gorm.DB
	dsn string
)

func init() {
	dsn = core_module.Config.GetString(dsnConfigKey)
	if dsn == "" {
		panic("database dsn is not set")
	}
}

func GetDefaultDB() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{
			Logger: NewSlogLogger(slog.Default().With("module", "gorm"),
				slogLoggerConfig{
					IngoreRunLog: core_module.Config.GetBool(loggerIgnoreRunLogConfigKey),
				}),
		})
		if err != nil {
			panic("failed to connect database")
		}
	}

	return db
}
