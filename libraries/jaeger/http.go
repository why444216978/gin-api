package jaeger

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

// ExtractHTTP is used to extract span context by HTTP middleware
func ExtractHTTP(ctx context.Context, req *http.Request, logID string) (context.Context, opentracing.Span, string) {
	var span opentracing.Span

	parentSpanContext, err := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if parentSpanContext == nil || err == opentracing.ErrSpanContextNotFound {
		span, ctx = opentracing.StartSpanFromContext(ctx, httpServerComponentPrefix+req.URL.Path)
	} else {
		span = Tracer.StartSpan(
			httpServerComponentPrefix+req.URL.Path,
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCServer,
		)
	}
	span.SetTag(string(ext.Component), httpServerComponentPrefix+req.URL.Path)
	span.SetTag(fieldLogID, logID)
	SetCommonTag(ctx, span)

	ctx = context.WithValue(opentracing.ContextWithSpan(ctx, span), parentSpanContextKey, span.Context())

	return ctx, span, GetSpanID(span)
}

// InjectHTTP is used to inject HTTP span
func InjectHTTP(ctx context.Context, req *http.Request, logID string) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, Tracer, httpClientComponentPrefix+req.URL.Path, ext.SpanKindRPCClient)
	defer span.Finish()
	span.SetTag(string(ext.Component), httpClientComponentPrefix+req.URL.Path)
	span.SetTag(fieldLogID, logID)
	SetCommonTag(ctx, span)

	err := Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-next-error", err.Error()))
	}

	return
}

func SetHTTPLog(span opentracing.Span, req, resp string) {
	span.LogFields(opentracing_log.Object(logFieldsRequest, string(req)))
	span.LogFields(opentracing_log.Object(logFieldsResponse, string(resp)))
}
