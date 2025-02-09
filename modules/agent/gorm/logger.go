package gorm_agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type slogLogger struct {
	logger *slog.Logger
	config slogLoggerConfig
}

type slogLoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	IngoreErrorLog            bool
	IgnoreSlowLog             bool
	IngoreRunLog              bool
}

func NewSlogLogger(logger *slog.Logger, config slogLoggerConfig) slogLogger {
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 100 * time.Millisecond
	}
	return slogLogger{
		logger: logger,
		config: config,
	}
}

func (l slogLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l slogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.InfoContext(ctx, msg, data...)
}

func (l slogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WarnContext(ctx, msg, data...)
}

func (l slogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ErrorContext(ctx, msg, data...)
}

func (l slogLogger) Debug(ctx context.Context, msg string, data ...interface{}) {
	l.logger.DebugContext(ctx, msg, data...)
}

func (l slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case !l.config.IngoreErrorLog &&
		err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		errLog := "Error SQL"
		if rows == -1 {
			l.Error(ctx, errLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql,
				"err", err)
		} else {
			l.Error(ctx, errLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql,
				"err", err)
		}
	case !l.config.IgnoreSlowLog &&
		elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("Slow SQL >= %v", l.config.SlowThreshold)
		if rows == -1 {
			l.Warn(ctx, slowLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql)
		} else {
			l.Warn(ctx, slowLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case !l.config.IngoreRunLog:
		sql, rows := fc()
		runLog := "Run SQL"
		if rows == -1 {
			l.Debug(ctx, runLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql)
		} else {
			l.Debug(ctx, runLog, "file", utils.FileWithLineNum(),
				"elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}
