package log

import (
	"context"
)

import (
	"github.com/glory-go/glory/config"
)

func init() {
	defaultLogger = NewLogger()
	defaultLogger.setup(config.GlobalServerConf.LogConfigs, config.GlobalServerConf.ServerName, config.GlobalServerConf.OrgName)
	defaultLogger.start()
}

var defaultLogger *Logger

func Debugf(template string, args ...interface{}) {
	defaultLogger.debugf(template, args...)
}

func CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	defaultLogger.ctxDebugf(ctx, template, args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.infof(template, args...)
}

func CtxInfof(ctx context.Context, template string, args ...interface{}) {
	defaultLogger.ctxInfof(ctx, template, args...)
}
func Warnf(template string, args ...interface{}) {
	defaultLogger.warnf(template, args...)
}

func CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	defaultLogger.ctxWarnf(ctx, template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.errorf(template, args...)
}

func CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	defaultLogger.ctxErrorf(ctx, template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultLogger.panicf(template, args...)
}

func CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	defaultLogger.ctxPanicf(ctx, template, args...)
}

func Debug(args ...interface{}) {
	defaultLogger.debug(args...)
}

func CtxDebug(ctx context.Context, args ...interface{}) {
	defaultLogger.ctxDebug(ctx, args...)
}

func Info(args ...interface{}) {
	defaultLogger.info(args...)
}

func CtxInfo(ctx context.Context, args ...interface{}) {
	defaultLogger.ctxInfo(ctx, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.warn(args...)
}

func CtxWarn(ctx context.Context, args ...interface{}) {
	defaultLogger.ctxWarn(ctx, args...)
}

func Error(args ...interface{}) {
	defaultLogger.error(args...)
}

func CtxError(ctx context.Context, args ...interface{}) {
	defaultLogger.ctxError(ctx, args...)
}

func Panic(args ...interface{}) {
	defaultLogger.panic(args...)
}

func CtxPanic(ctx context.Context, args ...interface{}) {
	defaultLogger.ctxPanic(ctx, args...)
}
