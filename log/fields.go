package log

import (
	"context"
	"sync"
)

import (
	"go.uber.org/zap"
)

var logFields sync.Map

func init() {
	logFields = sync.Map{}
}

type loaderFunc func(context.Context) interface{}

func WithField(key string, loader loaderFunc) {
	_, exist := logFields.LoadOrStore(key, loader)
	if exist {
		Warnf("field with name %v already been stored before", key)
	}
}

func With(ctx context.Context, l *zap.SugaredLogger) *zap.SugaredLogger {
	logFields.Range(func(key, value interface{}) bool {
		if value == nil {
			Warnf("loader of key %v is nil", key)
			return true
		}
		loader, ok := value.(loaderFunc)
		if !ok {
			Warnf("invalid type of key %v, it is %T", key, value)
			return true
		}
		value = loader(ctx)
		if value == nil {
			Infof("key %v got nil value", key)
			return true
		}
		l = l.With(key, value)
		return true
	})
	return l
}
