package jaeger

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := opentracing.GlobalTracer()

		opentracingSpanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if opentracingSpanContext == nil || err != nil {
			span := tracer.StartSpan(c.Request.URL.Path)
			defer span.Finish()
			opentracingSpanContext = span.Context()

			SetTag(c, span, opentracingSpanContext)
		}
		c.Set(FIELD_TRACER, tracer)
		c.Set(FIELD_SPAN_CONTEXT, opentracingSpanContext)

		c.Next()
	}
}
