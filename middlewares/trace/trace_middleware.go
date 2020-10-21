package trace

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

const defaultComponentName = "net/http"
const JaegerOpen = 1
const AppName = "gin-api"
const JaegerHostPort = "127.0.0.1:6831"

type Options struct {
	tracer        opentracing.Tracer
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

type OptionsFunc func(*Options)

func OperationNameFunc(f func(r *http.Request) string) OptionsFunc {
	return func(options *Options) {
		options.opNameFunc = f
	}
}

func WithTracer(tracer opentracing.Tracer) OptionsFunc {
	return func(options *Options) {
		options.tracer = tracer
	}
}

func WithComponentName(componentName string) OptionsFunc {
	return func(options *Options) {
		options.componentName = componentName
	}
}

func WithSpanObserver(f func(span opentracing.Span, r *http.Request)) OptionsFunc {
	return func(options *Options) {
		options.spanObserver = f
	}
}

func WithURLTagFunc(f func(u *url.URL) string) OptionsFunc {
	return func(options *Options) {
		options.urlTagFunc = f
	}
}

//OpenTracing 链路追踪中间件
//实现了[opentracing](https://opentracing.io)协议
//tracer默认使用jaeger.Tracer,如需修改,可用中间件WithTracer
func OpenTracing(serviceName string, options ...OptionsFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if JaegerOpen == 1 {

			var parentSpan opentracing.Span

			tracer, closer := NewJaegerTracer(AppName, JaegerHostPort)
			defer closer.Close()

			spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err != nil {
				parentSpan = tracer.StartSpan(c.Request.URL.Path)
				defer parentSpan.Finish()
			} else {
				parentSpan = opentracing.StartSpan(
					c.Request.URL.Path,
					opentracing.ChildOf(spCtx),
					opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
					ext.SpanKindRPCServer,
				)
				defer parentSpan.Finish()
			}
			c.Set("Tracer", tracer)
			c.Set("ParentSpanContext", parentSpan.Context())
		}
		c.Next()
	}
	// opts := Options{
	// 	opNameFunc: func(r *http.Request) string {
	// 		return r.Proto + " " + r.Method
	// 	},
	// 	spanObserver: func(span opentracing.Span, r *http.Request) {},
	// 	urlTagFunc: func(u *url.URL) string {
	// 		return u.String()
	// 	},
	// }
	// for _, opt := range options {
	// 	opt(&opts)
	// }

	// if opts.tracer == nil {
	// 	opts.tracer = opentracing.GlobalTracer()
	// }

	// return func(c *gin.Context) {
	// 	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	// 	spanContext, _ := opts.tracer.Extract(opentracing.HTTPHeaders, carrier)
	// 	op := opts.opNameFunc(c.Request)
	// 	sp := opts.tracer.StartSpan(op, opentracing.ChildOf(spanContext), opentracing.StartTime(time.Now()))
	// 	defer sp.Finish()
	// 	ext.HTTPMethod.Set(sp, c.Request.Method)
	// 	ext.HTTPUrl.Set(sp, opts.urlTagFunc(c.Request.URL))
	// 	opts.spanObserver(sp, c.Request)

	// 	componentName := opts.componentName
	// 	if componentName == "" {
	// 		componentName = defaultComponentName
	// 	}

	// 	ext.Component.Set(sp, componentName)
	// 	c.Request = c.Request.WithContext(
	// 		opentracing.ContextWithSpan(c.Request.Context(), sp))
	// 	//trace info 注入resp.Header
	// 	sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders,
	// 		opentracing.HTTPHeadersCarrier(c.Writer.Header()))
	// 	c.Next()
	// 	ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
	// 	sp.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
	// }
}

func NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer) {

	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: serviceName,
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
