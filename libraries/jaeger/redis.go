package jaeger

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

func InjectRedis(ctx context.Context, header http.Header, operationName, args string) (span opentracing.Span, err error) {
	parentSpanContext, ok := getInjectParent(ctx)
	if !ok {
		return
	}

	span = opentracing.StartSpan(
		operationName,
		opentracing.ChildOf(parentSpanContext),
		opentracing.Tag{Key: "args", Value: args},
		ext.SpanKindRPCClient,
	)
	setRedisTag(ctx, span)
	err = Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-current-error", err.Error()))
	}

	return
}

func setRedisTag(ctx context.Context, span opentracing.Span) {
	setTag(ctx, span)
	span.SetTag(string(ext.Component), operationTypeRedis)
}
