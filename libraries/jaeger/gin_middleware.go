package jaeger

import (
	"context"
	"gin-api/libraries/logging"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func GinOpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var span opentracing.Span
		ctx := c.Request.Context()

		parentSpanContext, err := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if parentSpanContext == nil || err != nil {
			span = Tracer.StartSpan(c.Request.URL.Path)
		} else {
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(parentSpanContext),
				ext.RPCServerOption(parentSpanContext),
				ext.SpanKindRPCClient,
			)
		}
		defer span.Finish()
		setHTTPTag(ctx, span)

		ctx = logging.AddTraceID(ctx, getTraceID(span))
		ctx = context.WithValue(opentracing.ContextWithSpan(ctx, span), parentSpanContextKey, span.Context())

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
