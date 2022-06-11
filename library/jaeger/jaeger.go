package jaeger

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	FieldLogID   = "Log-Id"
	FieldTraceID = "Trace-Id"
	FieldSpanID  = "Span-Id"
)

const (
	LogFieldsRequest  = "request"
	LogFieldsResponse = "response"
	LogFieldsArgs     = "args"
)

var Tracer opentracing.Tracer

type Config struct {
	Host string
	Port string
}

func NewJaegerTracer(connCfg *Config, serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", // 固定采样
			Param: 1,       // 1=全采样、0=不采样
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: connCfg.Host + ":" + connCfg.Port,
		},

		ServiceName: serviceName,
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	Tracer = tracer
	return tracer, closer, nil
}

func SetError(span opentracing.Span, err error) {
	span.LogFields(opentracingLog.Error(err))
}

func SetRequest(span opentracing.Span, request string) {
	span.LogFields(opentracingLog.String(LogFieldsRequest, request))
}

func SetResponse(span opentracing.Span, resp string) {
	span.LogFields(opentracingLog.String(LogFieldsResponse, resp))
}

func SetCommonTag(ctx context.Context, span opentracing.Span) {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	span.SetTag(FieldTraceID, jaegerSpanContext.TraceID().String())
	span.SetTag(FieldSpanID, jaegerSpanContext.SpanID().String())
}

func GetTraceID(span opentracing.Span) string {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	return jaegerSpanContext.TraceID().String()
}

func GetSpanID(span opentracing.Span) string {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	return jaegerSpanContext.SpanID().String()
}

func spanContextToJaegerContext(spanContext opentracing.SpanContext) jaeger.SpanContext {
	if sc, ok := spanContext.(jaeger.SpanContext); ok {
		return sc
	} else {
		return jaeger.SpanContext{}
	}
}
