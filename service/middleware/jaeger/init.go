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
	)
}
