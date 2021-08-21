package opentracing

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"gin-api/global"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/response"

	"github.com/gin-gonic/gin"
)

func Rpc(c *gin.Context) {
	sendUrl := fmt.Sprintf("http://localhost:%d/test/rpc1?logid=%s", global.Global.AppPort, logging.ValueLogID(c))

	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := jaeger.JaegerSend(c.Request.Context(), http.MethodPost, sendUrl, header, bytes.NewBufferString(`{"a":"a"}`), time.Second)
	if err != nil {
		fmt.Println(ret)
		fmt.Println(err)
		return
	}

	response.Response(c, response.CodeSuccess, ret, "")
}

func Rpc1(c *gin.Context) {
	sendUrl := fmt.Sprintf("http://localhost:%d/test/conn?logid=%s", global.Global.AppPort, logging.ValueLogID(c))

	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := jaeger.JaegerSend(c.Request.Context(), http.MethodPost, sendUrl, header, nil, time.Second)
	if err != nil {
		fmt.Println(ret)
		fmt.Println(err)
		return
	}

	response.Response(c, response.CodeSuccess, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
