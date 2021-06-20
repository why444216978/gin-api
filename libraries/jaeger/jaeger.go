package jaeger

import (
	"gin-api/global"
	"gin-api/libraries/logging"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	FIELD_LOG_ID       = "Log-Id"
	FIELD_TRACE_ID     = "Trace-Id"
	FIELD_SPAN_ID      = "Span-Id"
	FIELD_TRACER       = "Tracer"
	FIELD_SPAN_CONTEXT = "Span"
)

const (
	OPERATION_TYPE_HTTP     = "HTTP"
	OPERATION_TYPE_RPC      = "RPC"
	OPERATION_TYPE_MYSQL    = "MySQL"
	OPERATION_TYPE_REDIS    = "Redis"
	OPERATION_TYPE_RabbitMQ = "RabbitMQ"
)

var (
	Tracer opentracing.Tracer
)

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

func InjectHTTP(c *gin.Context, header http.Header, operationName, operationType string) (span opentracing.Span, err error) {
	parentSpanContext, ok := getInjectParent(c)
	if !ok {
		return
	}
	err = Tracer.Inject(parentSpanContext, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-next-error", err.Error()))
	}

	return
}

func InjectRedis(c *gin.Context, header http.Header, operationName, args string) (span opentracing.Span, err error) {
	parentSpanContext, ok := getInjectParent(c)
	if !ok {
		return
	}

	span = opentracing.StartSpan(
		operationName,
		opentracing.ChildOf(parentSpanContext),
		opentracing.Tag{Key: string(ext.Component), Value: OPERATION_TYPE_REDIS},
		opentracing.Tag{Key: "args", Value: args},
		ext.SpanKindRPCClient,
	)
	SetTag(c, span, parentSpanContext)
	err = Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-current-error", err.Error()))
	}

	return
}

func getInjectParent(c *gin.Context) (spanContext opentracing.SpanContext, ok bool) {
	var _spanContext interface{}

	_spanContext, ok = c.Get(FIELD_SPAN_CONTEXT)
	if !ok {
		return
	}
	spanContext, ok = _spanContext.(opentracing.SpanContext)
	if !ok {
		return
	}

	return
}

func SetTag(c *gin.Context, span opentracing.Span, spanContext opentracing.SpanContext) {
	jaegerSpanContext := spanContextToJaegerContext(spanContext)
	span.SetTag(FIELD_TRACE_ID, jaegerSpanContext.TraceID().String())
	span.SetTag(FIELD_SPAN_ID, jaegerSpanContext.SpanID().String())
	span.SetTag(FIELD_LOG_ID, logging.ValueLogID(c))
}

func spanContextToJaegerContext(spanContext opentracing.SpanContext) jaeger.SpanContext {
	if sc, ok := spanContext.(jaeger.SpanContext); ok {
		return sc
	} else {
		return jaeger.SpanContext{}
	}
}
