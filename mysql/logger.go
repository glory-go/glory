package mysql

import (
	"context"
	syslog "log"
	"os"
	"time"
)

import (
	"gorm.io/gorm/logger"
)

import (
	"github.com/glory-go/glory/log"
)

type GormLogger struct {
	level logger.LogLevel

	defaultLogger logger.Interface
}

func NewGormLogger() logger.Interface {
	return &GormLogger{
		defaultLogger: logger.New(
			syslog.New(os.Stdout, "\r\n", syslog.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      false,
			},
		),
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.defaultLogger = l.defaultLogger.LogMode(level)
	l.level = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.defaultLogger.Info(ctx, msg, data...)
	if l.level >= logger.Info {
		log.CtxInfof(ctx, msg, data)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.defaultLogger.Warn(ctx, msg, data...)
	if l.level >= logger.Warn {
		log.CtxWarnf(ctx, msg, data)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.defaultLogger.Error(ctx, msg, data...)
	if l.level >= logger.Error {
		log.CtxErrorf(ctx, msg, data)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.defaultLogger.Trace(ctx, begin, fc, err)
}
