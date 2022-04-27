package jaeger

import (
	"io"
)

import (
	"github.com/opentracing/opentracing-go"

	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/grmanager"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/tools"
)

const (
	AliyunJaegerConfigKey = "aliyun_jaeger"
)

var (
	tracer opentracing.Tracer
)

type aliyunJaegerConfig struct {
	ConfigSource string `mapstructure:"config_source"`
	Endpoint     string `mapstructure:"endpoint"`
}

func init() {
	viper := config.GetViperConfig()
	jaegerConfig := &aliyunJaegerConfig{}
	if err := viper.UnmarshalKey(AliyunJaegerConfigKey, jaegerConfig); err != nil {
		log.Warnf("jager fail to parse config with err %v", err)
		return
	}
	if err := tools.ReadFromEnvIfNeed(jaegerConfig); err != nil {
		log.Warnf("jager fail to read config in env with err %v", err)
		return
	}
	sender := transport.NewHTTPTransport(
		jaegerConfig.Endpoint,
	)
	var closer io.Closer
	tracer, closer = jaeger.NewTracer(config.GlobalServerConf.GetAppKey(),
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender),
		jaeger.TracerOptions.Logger(&jaegerLogger{}),
	)
	grmanager.RegisterCloser(closer)
}

type jaegerLogger struct{}

func (l *jaegerLogger) Error(msg string) {
	log.Error(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}
