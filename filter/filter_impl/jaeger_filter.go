package filter_impl

import (
	"context"
	"io"
	"strings"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/filter"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaegerlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	plugin.SetFilterFactory("jaeger", NewJaegerFilter)
}

type JaegerFilter struct {
	tracer            opentracing.Tracer
	aliyunUploadToken string
	next              filter.GRPCFilter
}

func NewJaegerFilter(filterConfig *config.FilterConfig) (filter.GRPCFilter, error) {
	jaegerFilter := &JaegerFilter{}

	conf, err := jaegerFilter.setup(filterConfig)
	if err != nil {
		log.Errorf("ERROR: fail to setup Jaeger:%v\n", err)
		return nil, err
	}
	jaegerFilter.tracer, _, err = jaegerFilter.newJaegerTracer(conf)
	if err != nil {
		log.Errorf("ERROR: fail to init Jaeger:%v\n", err)
		return nil, err
	}

	return jaegerFilter, nil
}

func (j *JaegerFilter) SetNext(filter filter.GRPCFilter) {
	j.next = filter
}

func (j *JaegerFilter) ServerHandle(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo) (resp interface{}, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	spanContext, err := j.tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		log.Errorf("JaegerError: extract from metadata err: %v", err)
	} else {
		span := j.tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCServer,
		)
		defer span.Finish()
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			log.Info("Tracing " + config.GlobalServerConf.ServerName + " trace-id:" + sc.TraceID().String() + " span-id:" + sc.SpanID().String())
		} else if !ok {
			log.Error("JaegerError: fail to get id")
		}

		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	if j.next == nil {
		log.Errorf("err = %v", filter.ErrNotSetNextFilter)
		return nil, filter.ErrNotSetNextFilter
	}
	return j.next.ServerHandle(ctx, req, info)
}

func (j *JaegerFilter) ClientHandle(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	span := j.tracer.StartSpan(
		method,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCClient,
	)
	defer span.Finish()
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		log.Debug("Tracing " + config.GlobalServerConf.ServerName + " trace-id:" + sc.TraceID().String() + " span-id:" + sc.SpanID().String())
	} else if !ok {
		log.Error("JaegerError: fail to get id")
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md.Copy()
	}

	mdWriter := MDReaderWriter{md}
	err := j.tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
	if err != nil {
		span.LogFields(jaegerlog.String("inject-error", err.Error()))
		log.Errorf("JaegerError: inject-error", err.Error())
	}

	newCtx := metadata.NewOutgoingContext(ctx, md)

	if j.next == nil {
		log.Errorf("err = %v", filter.ErrNotSetNextFilter)
		return filter.ErrNotSetNextFilter
	}
	err = j.next.ClientHandle(newCtx, method, req, reply, cc, invoker, opts...)
	if err != nil {
		span.LogFields(jaegerlog.String("call-error", err.Error()))
		log.Errorf("JaegerError: call-error", err.Error())
	}
	return err

}

func (j *JaegerFilter) setup(conf *config.FilterConfig) (*jaegercfg.Configuration, error) {
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  conf.SamplerType,
			Param: conf.SamplerParam, // todo: what's meaning?
		},

		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: conf.Address,
		},

		ServiceName: config.GlobalServerConf.ServerName,
	}

	j.aliyunUploadToken = conf.AliyunToken1 + "_" + conf.AliyunToken2
	return cfg, nil
}

// 根据config新建tracer
func (j *JaegerFilter) newJaegerTracer(jcfg *jaegercfg.Configuration) (tracer opentracing.Tracer, closer io.Closer, err error) {
	sender := transport.NewHTTPTransport(
		"http://tracing-analysis-dc-hz.aliyuncs.com/adapt_" + j.aliyunUploadToken + "/api/traces",
	)
	tracer, closer, err = jcfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Sampler(jaeger.NewConstSampler(true)),
		jaegercfg.Reporter(jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger))),
	)
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}

//MDReaderWriter metadata Reader and Writer
type MDReaderWriter struct {
	metadata.MD
}

// ForeachKey implements ForeachKey of opentracing.TextMapReader
func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter
func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}
