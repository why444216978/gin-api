package jaeger

import (
	"context"
	"errors"
	"fmt"
	"gin-api/libraries/logging"
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

func ExtractHTTP(ctx context.Context, req *http.Request) (context.Context, opentracing.Span) {
	var span opentracing.Span

	parentSpanContext, err := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if parentSpanContext == nil || err != nil {
		span = Tracer.StartSpan(req.URL.Path)
	} else {
		span = opentracing.StartSpan(
			req.URL.Path,
			opentracing.ChildOf(parentSpanContext),
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCClient,
		)
	}
	span.SetTag(string(ext.Component), operationTypeHTTP)

	SetCommonTag(ctx, span)

	ctx = logging.AddTraceID(ctx, GetTraceID(span))
	ctx = context.WithValue(opentracing.ContextWithSpan(ctx, span), parentSpanContextKey, span.Context())

	return ctx, span
}

// JaegerSend 发送Jaeger请求
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
	opentracingSpan, _ := injectHTTP(ctx, req.Header, req.URL.Path, operationTypeHTTP)
	if opentracingSpan != nil {
		defer opentracingSpan.Finish()
	}

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

func injectHTTP(ctx context.Context, header http.Header, operationName, operationType string) (span opentracing.Span, err error) {
	parentSpanContext, ok := getInjectParent(ctx)
	if !ok {
		return
	}
	err = Tracer.Inject(parentSpanContext, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span.LogFields(opentracing_log.String("inject-next-error", err.Error()))
	}

	return
}

func SetHTTPLog(span opentracing.Span, req, resp string) {
	span.LogFields(opentracing_log.Object(logFieldsRequest, string(req)))
	span.LogFields(opentracing_log.Object(logFieldsResponse, string(resp)))
}
