package jaeger

import (
	"errors"
	"fmt"
	"gin-api/libraries/logging"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

type Response struct {
	HTTPCode int
	Response string
}

// JaegerSend 发送Jaeger请求
func JaegerSend(c *gin.Context, method, url string, header map[string]string, body io.Reader, timeout time.Duration) (ret Response, err error) {
	var req *http.Request

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	//构建req
	req, err = http.NewRequestWithContext(c, method, url, body)
	if err != nil {
		return
	}

	//设置请求header
	req.Header.Add(logging.LOG_FIELD, logging.ValueLogID(c))
	for k, v := range header {
		req.Header.Add(k, v)
	}

	//注入
	tracer, ok1 := c.Get(FIELD_TRACER)
	parentSpanContext, ok2 := c.Get(FIELD_SPAN)
	if ok1 && ok2 {
		spParent := parentSpanContext.(opentracing.Span)
		span := opentracing.StartSpan(
			req.URL.Path,
			opentracing.ChildOf(spParent.Context()),
		)
		defer span.Finish()

		sp := spanContextToJaegerContext(span.Context())
		span.SetTag(FIELD_TRACE_ID, sp.TraceID().String())
		span.SetTag(FIELD_SPAN_ID, sp.SpanID().String())
		span.SetTag(FIELD_LOG_ID, logging.ValueLogID(c))

		err = tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			span.LogFields(opentracing_log.String("inject-error", err.Error()))
			err = nil
		}
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
