package trace

import (
	"context"
	"gin-frame/libraries/log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const defaultComponentName = "net/http"

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
	opts := Options{
		opNameFunc: func(r *http.Request) string {
			return r.Proto + " " + r.Method
		},
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.tracer == nil {
		opts.tracer = opentracing.GlobalTracer()
	}

	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		spanContext, _ := opts.tracer.Extract(opentracing.HTTPHeaders, carrier)
		op := opts.opNameFunc(c.Request)
		sp := opts.tracer.StartSpan(op, opentracing.ChildOf(spanContext), opentracing.StartTime(time.Now()))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, opts.urlTagFunc(c.Request.URL))
		opts.spanObserver(sp, c.Request)

		componentName := opts.componentName
		if componentName == "" {
			componentName = defaultComponentName
		}

		ext.Component.Set(sp, componentName)
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), sp))
		//trace info 注入resp.Header
		sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Writer.Header()))
		c.Next()
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
		sp.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
	}
}

//rpc调用时，trace注入header，用于分析当前rpc调用链的各节点耗时
func InjectTrace(ctx context.Context, logFormater *log.LogFormat, req *http.Request) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		panic("nil span")
	}

	//trace info 注入req.Header
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	req.Header.Set("x-hop", logFormater.XHop.Hex())
}
