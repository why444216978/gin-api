package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/libraries/registry"
)

type Response struct {
	HTTPCode int
	Response string
}

// Send is send HTTP request
func Send(ctx context.Context, serviceName, method, uri string, header map[string]string, body io.Reader, timeout time.Duration) (ret *Response, err error) {
	ser := registry.Services["gin-api"]
	if ser == nil {
		return nil, errors.New("service is nil")
	}

	node := ser.GetServices()
	if len(node) <= 0 {
		return nil, errors.New("service node empty")
	}

	var req *http.Request

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	//构建req
	req, err = http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://%s:%d%s", node[0].Host, node[0].Port, uri), body)
	if err != nil {
		return
	}

	//设置请求header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if ctx.Err() != nil {
		return nil, err
	}

	//注入Jaeger
	logID := req.Header.Get(logging.LogHeader)
	jaeger.InjectHTTP(ctx, req, logID)

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

	ret = &Response{
		HTTPCode: resp.StatusCode,
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http code is %d", resp.StatusCode))
		return
	}

	if b != nil {
		ret.Response = string(b)
	}

	return
}
