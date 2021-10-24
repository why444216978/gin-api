package jaeger

import (
	"context"
	"net/http"

	"github.com/why444216978/gin-api/libraries/jaeger"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

const (
	httpClientComponentPrefix = "HTTP-Client-"
	httpServerComponentPrefix = "HTTP-Server-"
)

// ExtractHTTP is used to extract span context by HTTP middleware
func ExtractHTTP(ctx context.Context, req *http.Request, logID string) (context.Context, opentracing.Span, string) {
	if jaeger.Tracer == nil {
		return ctx, nil, ""
	}

	var span opentracing.Span

	parentSpanContext, err := jaeger.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if parentSpanContext == nil || err == opentracing.ErrSpanContextNotFound {
		span, ctx = opentracing.StartSpanFromContext(ctx, httpServerComponentPrefix+req.URL.Path)
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

	ctx = context.WithValue(opentracing.ContextWithSpan(ctx, span), jaeger.ParentSpanContextKey, span.Context())

	return ctx, span, jaeger.GetSpanID(span)
}

// InjectHTTP is used to inject HTTP span
func InjectHTTP(ctx context.Context, req *http.Request, logID string) {
	if jaeger.Tracer == nil {
		return
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, jaeger.Tracer, httpClientComponentPrefix+req.URL.Path, ext.SpanKindRPCClient)
	defer span.Finish()
	span.SetTag(string(ext.Component), httpClientComponentPrefix+req.URL.Path)
	span.SetTag(jaeger.FieldLogID, logID)
	jaeger.SetCommonTag(ctx, span)

	err := jaeger.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-next-error", err.Error()))
	}

	return
}

func SetHTTPLog(span opentracing.Span, req, resp string) {
	span.LogFields(opentracing_log.Object(jaeger.LogFieldsRequest, string(req)))
	span.LogFields(opentracing_log.Object(jaeger.LogFieldsResponse, string(resp)))
}
