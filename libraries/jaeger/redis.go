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
		opentracing.Tag{Key: string(ext.Component), Value: operationTypeRedis},
		ext.SpanKindRPCClient,
	)
	SetCommonTag(ctx, span)

	span.LogFields(opentracing_log.String(logFieldsArgs, args))

	err = Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span.LogFields(opentracing_log.Error(err))
	}

	return
}
