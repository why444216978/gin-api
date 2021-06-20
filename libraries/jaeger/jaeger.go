package jaeger

import (
	"context"
	"gin-api/global"
	"gin-api/libraries/logging"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

type Config struct {
	Host string
	Port string
}

func NewJaegerTracer(connCfg Config) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: connCfg.Host + ":" + connCfg.Port,
		},

		ServiceName: global.Global.AppName,
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	Tracer = tracer
	return tracer, closer, nil
}

func setTag(ctx context.Context, span opentracing.Span) {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	span.SetTag(fieldTraceID, jaegerSpanContext.TraceID().String())
	span.SetTag(fieldSpanID, jaegerSpanContext.SpanID().String())
	span.SetTag(fieldLogID, logging.ValueLogID(ctx))
}

func getInjectParent(ctx context.Context) (spanContext opentracing.SpanContext, ok bool) {
	var _spanContext interface{}

	_spanContext = ctx.Value(parentSpanContextKey)
	spanContext, ok = _spanContext.(opentracing.SpanContext)
	if !ok {
		return
	}

	return
}

func spanContextToJaegerContext(spanContext opentracing.SpanContext) jaeger.SpanContext {
	if sc, ok := spanContext.(jaeger.SpanContext); ok {
		return sc
	} else {
		return jaeger.SpanContext{}
	}
}
