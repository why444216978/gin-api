package http

import (
	"context"
	"errors"
	"fmt"
	"gin-api/libraries/jaeger"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	HTTPCode int
	Response string
}

// Send is send HTTP request
func Send(ctx context.Context, method, url string, header map[string]string, body io.Reader, timeout time.Duration) (ret Response, err error) {
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
	jaeger.InjectHTTP(ctx, req)

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
