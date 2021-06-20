package jaeger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			span opentracing.Span
			ctx  context.Context
		)
		parentSpanContext, err := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if parentSpanContext == nil || err != nil {
			// span = Tracer.StartSpan(c.Request.URL.Path)
			span, ctx = opentracing.StartSpanFromContext(c.Request.Context(), c.Request.URL.Path)
		} else {
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				// opentracing.ChildOf(parentSpanContext),
				ext.RPCServerOption(parentSpanContext),
				opentracing.Tag{Key: string(ext.Component), Value: OPERATION_TYPE_HTTP},
				ext.SpanKindRPCClient,
			)
		}
		defer span.Finish()
		SetTag(c, span, span.Context())

		c.Set(FIELD_TRACER, Tracer)
		c.Set(FIELD_SPAN_CONTEXT, span.Context())

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
