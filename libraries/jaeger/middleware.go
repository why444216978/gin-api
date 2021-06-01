package jaeger

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := opentracing.GlobalTracer()

		var span opentracing.Span
		parentSpanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if parentSpanContext == nil || err != nil {
			span = tracer.StartSpan(c.Request.URL.Path)
		} else {
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(parentSpanContext),
				opentracing.Tag{Key: string(ext.Component), Value: OPERATION_TYPE_HTTP},
				ext.SpanKindRPCClient,
			)
		}
		defer span.Finish()
		SetTag(c, span, span.Context())

		c.Set(FIELD_TRACER, tracer)
		c.Set(FIELD_SPAN_CONTEXT, span.Context())

		c.Next()
	}
}
