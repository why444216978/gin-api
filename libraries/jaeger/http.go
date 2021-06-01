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

	//注入Jaeger
	opentracingSpan, _ := InjectHTTP(c, req.Header, c.Request.URL.Path, OPERATION_TYPE_HTTP)
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
