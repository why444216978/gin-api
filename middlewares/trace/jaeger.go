package trace

import (
	"fmt"
	"gin-api/libraries/jaeger"
	"gin-api/resource"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	// "github.com/opentracing/opentracing-go/log"
)

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := resource.Tracer

		var sp opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			fmt.Println("opentracing:start")
			sp = tracer.StartSpan(c.Request.URL.Path)
		} else {
			fmt.Println("opentracing:extract")
			sp = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				// opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				// ext.SpanKindRPCServer,
			)
		}
		defer sp.Finish()

		span := jaeger.SpanContextToJaegerContext(sp.Context())
		sp.SetTag(jaeger.FIELD_TRACE_ID, span.TraceID().String())
		sp.SetTag(jaeger.FIELD_SPAN_ID, span.SpanID().String())

		c.Next()
	}
}
