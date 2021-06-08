package log

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/glory-go/glory/config"
	"github.com/natefinch/lumberjack"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	rawLogConfig *config.LogConfig
	logLevel     zapcore.Level
	logType      string
	filePath     string
	elasticAddr  string
	serviceName  string
	orgName      string
}

//Logger 一个logger对应多个实例化zapcore，调用Logger的方法将作用于所有core，从而实现日志的多重记录
type Logger struct {
	config   map[string]*LoggerConfig
	logger   []*zap.SugaredLogger
	enconfig zapcore.EncoderConfig
}

func NewLogger() *Logger {
	return &Logger{
		config: make(map[string]*LoggerConfig),
	}
}

func (l *Logger) setup(logConfigMap map[string]*config.LogConfig, serviceName, orgName string) {
	for k, v := range logConfigMap {
		l.config[k] = &LoggerConfig{
			logType:      v.LogType,
			filePath:     v.FilePath,
			elasticAddr:  v.ElasticAddr,
			serviceName:  serviceName,
			orgName:      orgName,
			rawLogConfig: v,
		}
		switch v.LogLevel {
		case "debug":
			l.config[k].logLevel = zap.DebugLevel
		case "info":
			l.config[k].logLevel = zap.InfoLevel
		case "warn":
			l.config[k].logLevel = zap.WarnLevel
		case "error":
			l.config[k].logLevel = zap.ErrorLevel
		case "panic":
			l.config[k].logLevel = zap.PanicLevel
		default:
			l.config[k].logLevel = zap.DebugLevel
		}
	}
}

func (l *Logger) start() {
	// 生成编码配置
	enconfig := zap.NewProductionEncoderConfig()
	enconfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, v := range l.config {
		var newLogger *zap.SugaredLogger
		switch v.logType {
		case "console":
			newLogger = getConsoleSugaredLogger(v.logLevel, enconfig)
		case "file":
			newLogger = getFileSugaredLogger(v.logLevel, v.filePath, enconfig)
		case "elastic":
			newLogger = getElasticSugaredLogger(v.logLevel, v.elasticAddr, enconfig, v.serviceName)
		case "sls":
			newLogger = getAliyunSLSLogger(v.logLevel, enconfig, v.rawLogConfig, v.serviceName, v.orgName)
		}
		l.logger = append(l.logger, newLogger)
	}
}

func getElasticSugaredLogger(level zapcore.Level, addr string, enconfig zapcore.EncoderConfig, serviceName string) *zap.SugaredLogger {
	enconfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoder := zapcore.NewConsoleEncoder(enconfig)

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > lvl {
			return false
		}
		return lvl == zapcore.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > lvl {
			return false
		}
		return lvl == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > lvl {
			return false
		}
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > lvl {
			return false
		}
		return lvl == zapcore.ErrorLevel
	})
	panicLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > lvl {
			return false
		}
		return lvl == zapcore.PanicLevel
	})

	es, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(addr),
	)
	if err != nil {
		log.Println("elastic init err = " + err.Error())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(ElasticDebugHook{es: es, serviceName: serviceName}), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(ElasticInfoHook{es: es, serviceName: serviceName}), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(ElasticWarnHook{es: es, serviceName: serviceName}), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(ElasticErrorHook{es: es, serviceName: serviceName}), errorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(ElasticPanicHook{es: es, serviceName: serviceName}), panicLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	return logger.Sugar()
}

func getConsoleSugaredLogger(level zapcore.Level, enconfig zapcore.EncoderConfig) *zap.SugaredLogger {
	w := zapcore.AddSync(&ConsoleHook{})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enconfig), //编码器配置
		w,                                   //打印到控制台和文件
		level,                               //日志等级
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	return logger.Sugar()
}

func getFileSugaredLogger(level zapcore.Level, filePath string, enconfig zapcore.EncoderConfig) *zap.SugaredLogger {
	hook := &lumberjack.Logger{
		Filename:   fmt.Sprintf(filePath), //filePath
		MaxSize:    500,                   // megabytes
		MaxBackups: 10000,
		MaxAge:     100000, //days
		Compress:   false,  // disabled by default
	}
	defer hook.Close()
	w := zapcore.AddSync(hook)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enconfig), //编码器配置
		w,                                   //打印到控制台和文件
		level,                               //日志等级
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	return logger.Sugar()
}

func getAliyunSLSLogger(level zapcore.Level, enconfig zapcore.EncoderConfig, config *config.LogConfig, serverName, orgName string) *zap.SugaredLogger {
	core := newAliyunSLSLoggerCore(
		zapcore.NewConsoleEncoder(enconfig),
		level,
		config,
		serverName,
		orgName,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	return logger.Sugar()
}

func getremoteSugaredLogger(level zapcore.Level) *zap.SugaredLogger {
	// todo 远程日志
	return nil
}

func getTraceIDFromUberCtx(ctx context.Context) ([]interface{}, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
		return []interface{}{}, false
	}
	traceIDFields := md.Get("uber-trace-id")
	if len(traceIDFields) == 0 {
		return []interface{}{}, false
	}
	idList := strings.Split(traceIDFields[0], ":")
	if len(idList) == 0 {
		return []interface{}{}, false
	}
	//todo use to extend ctx field as user wanted
	return []interface{}{"uber-trace-id", idList[0]}, true
}

func (l *Logger) debugf(template string, arg ...interface{}) {
	for _, v := range l.logger {
		v.Debugf(template, arg...)
	}
}

func (l *Logger) ctxDebugf(ctx context.Context, template string, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.debugf(template, arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Debugf(template, arg...)
	}
}

func (l *Logger) infof(template string, arg ...interface{}) {
	for _, v := range l.logger {
		v.Infof(template, arg...)
	}
}

func (l *Logger) ctxInfof(ctx context.Context, template string, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.infof(template, arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Infof(template, arg...)
	}
}

func (l *Logger) warnf(template string, arg ...interface{}) {
	for _, v := range l.logger {
		v.Warnf(template, arg...)
	}
}

func (l *Logger) ctxWarnf(ctx context.Context, template string, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.warnf(template, arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Warnf(template, arg...)
	}
}

func (l *Logger) errorf(template string, arg ...interface{}) {
	arg = append(arg, zap.StackSkip("stack", 2).String)
	template += "\n%s\n"
	for _, v := range l.logger {
		v.Errorf(template, arg...)
	}
}

func (l *Logger) ctxErrorf(ctx context.Context, template string, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	template += "\n%s\n"
	arg = append(arg, zap.StackSkip("stack", 2).String)
	if !ok {
		l.errorf(template, arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Errorf(template, arg...)
	}
}

func (l *Logger) panicf(template string, arg ...interface{}) {
	for _, v := range l.logger {
		v.Panicf(template, arg...)
	}
}
func (l *Logger) ctxPanicf(ctx context.Context, template string, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.panicf(template, arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Panicf(template, arg...)
	}
}

func (l *Logger) debug(arg ...interface{}) {
	for _, v := range l.logger {
		v.Debug(arg...)
	}
}

func (l *Logger) ctxDebug(ctx context.Context, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.debug(arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Debug(arg...)
	}
}

func (l *Logger) info(arg ...interface{}) {
	for _, v := range l.logger {
		v.Info(arg...)
	}
}

func (l *Logger) ctxInfo(ctx context.Context, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.info(arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Info(arg...)
	}
}

func (l *Logger) warn(arg ...interface{}) {
	for _, v := range l.logger {
		v.Warn(arg...)
	}
}

func (l *Logger) ctxWarn(ctx context.Context, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.warn(arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Warn(arg...)
	}
}

func (l *Logger) error(arg ...interface{}) {
	arg = append(arg, zap.StackSkip("stack", 2).String)

	for _, v := range l.logger {
		v.Error(arg...)
	}
}
func (l *Logger) ctxError(ctx context.Context, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.error(arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Error(arg...)
	}
}

func (l *Logger) panic(arg ...interface{}) {
	for _, v := range l.logger {
		v.With(zap.StackSkip("stack", 2))
		v.Panic(arg...)
	}
}

func (l *Logger) ctxPanic(ctx context.Context, arg ...interface{}) {
	traceID, ok := getTraceIDFromUberCtx(ctx)
	if !ok {
		l.panic(arg...)
		return
	}
	for _, v := range l.logger {
		v.With(traceID).Panic(arg...)
	}
}
