package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/why444216978/go-util/assert"

	"github.com/why444216978/gin-api/library/jaeger"
)

const (
	httpClientComponentPrefix = "HTTP-Client-"
	httpServerComponentPrefix = "HTTP-Server-"
)

var ErrTracerNil = errors.New("Tracer is nil")

// ExtractHTTP is used to extract span context by HTTP middleware
func ExtractHTTP(ctx context.Context, req *http.Request, logID string) (context.Context, opentracing.Span, string) {
	if assert.IsNil(jaeger.Tracer) {
		return ctx, nil, ""
	}

	var span opentracing.Span

	parentSpanContext, err := jaeger.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if assert.IsNil(parentSpanContext) || err == opentracing.ErrSpanContextNotFound {
		span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, jaeger.Tracer, httpServerComponentPrefix+req.URL.Path, ext.SpanKindRPCServer)
	} else {
		span = jaeger.Tracer.StartSpan(
			httpServerComponentPrefix+req.URL.Path,
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCServer,
		)
	}
	span.SetTag(string(ext.Component), httpServerComponentPrefix+req.URL.Path)
	span.SetTag(jaeger.FieldLogID, logID)
	jaeger.SetCommonTag(ctx, span)

	ctx = opentracing.ContextWithSpan(ctx, span)

	return ctx, span, jaeger.GetTraceID(span)
}

// InjectHTTP is used to inject HTTP span
func InjectHTTP(ctx context.Context, req *http.Request, logID string) error {
	if assert.IsNil(jaeger.Tracer) {
		return ErrTracerNil
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, jaeger.Tracer, httpClientComponentPrefix+req.URL.Path, ext.SpanKindRPCClient)
	defer span.Finish()
	span.SetTag(string(ext.Component), httpClientComponentPrefix+req.URL.Path)
	span.SetTag(jaeger.FieldLogID, logID)
	jaeger.SetCommonTag(ctx, span)

	return jaeger.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}

func SetHTTPLog(span opentracing.Span, req, resp string) {
	if assert.IsNil(span) {
		return
	}
	jaeger.SetRequest(span, req)
	jaeger.SetResponse(span, resp)
}
