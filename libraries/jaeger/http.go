package jaeger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

type Response struct {
	HTTPCode int
	Response string
}

// ExtractHTTP is used to extract span context by HTTP middleware
func ExtractHTTP(ctx context.Context, req *http.Request, logID string) (context.Context, opentracing.Span, string) {
	var span opentracing.Span

	parentSpanContext, err := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if parentSpanContext == nil || err == opentracing.ErrSpanContextNotFound {
		span, ctx = opentracing.StartSpanFromContext(ctx, serverTagPrefix+req.URL.Path)
	} else {
		span = Tracer.StartSpan(
			serverTagPrefix+req.URL.Path,
			ext.RPCServerOption(parentSpanContext),
			opentracing.Tag{Key: string(ext.Component), Value: operationTypeHTTP},
			opentracing.Tag{Key: fieldLogID, Value: logID},
			ext.SpanKindRPCServer,
		)
	}
	span.SetTag(string(ext.Component), serverTagPrefix)
	SetCommonTag(ctx, span)

	ctx = context.WithValue(opentracing.ContextWithSpan(ctx, span), parentSpanContextKey, span.Context())

	return ctx, span, GetSpanID(span)
}

// JaegerSend send HTTP and inject span
func JaegerSend(ctx context.Context, method, url string, header map[string]string, body io.Reader, timeout time.Duration) (ret Response, err error) {
	var req *http.Request

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	//构建req
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	//设置请求header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	//注入Jaeger
	InjectHTTP(ctx, req)

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ret.HTTPCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http code is %d", resp.StatusCode))
		return
	}

	if b != nil {
		ret.Response = string(b)
	}

	return
}

// InjectHTTP is used to inject HTTP span
func InjectHTTP(ctx context.Context, req *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, Tracer, clientTagPrefix+req.URL.Path, ext.SpanKindRPCClient)
	defer span.Finish()
	span.SetTag(string(ext.Component), operationTypeHTTP)

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
