package jaeger

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

var (
	tracer opentracing.Tracer
)

type aliyunJaegerConfig struct {
	Endpoint string
}

func init() {
	viper := config.GetViperConfig()
	jaegerConfig := &aliyunJaegerConfig{}
	if err := viper.Unmarshal(jaegerConfig); err != nil {
		log.Warnf("jager fail to parse config with err %v", err)
		return
	}
	sender := transport.NewHTTPTransport(
		jaegerConfig.Endpoint,
	)
	tracer, _ = jaeger.NewTracer(config.GlobalServerConf.GetAppKey(),
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender),
		jaeger.TracerOptions.Logger(&jaegerLogger{}),
	)
}

type jaegerLogger struct{}

func (l *jaegerLogger) Error(msg string) {
	log.Error(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}
