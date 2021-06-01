package opentracing

import (
	"fmt"
	"net/http"
	"time"

	"gin-api/app_const"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/response"

	"github.com/gin-gonic/gin"
)

func Rpc(c *gin.Context) {
	sp, _ := jaeger.Inject(c, c.Request.Header, "select", jaeger.OPERATION_TYPE_MYSQL)
	if sp != nil {
		defer sp.Finish()
	}

	sendUrl := fmt.Sprintf("http://localhost:%d/test/rpc1?logid=%s", app_const.SERVICE_PORT, logging.ValueLogID(c))

	ret, err := jaeger.JaegerSend(c, http.MethodPost, sendUrl, nil, nil, time.Second)
	ret, err = jaeger.JaegerSend(c, http.MethodPost, sendUrl, nil, nil, time.Second)
	if err != nil {
		fmt.Println(err)
	}

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Rpc1(c *gin.Context) {
	sendUrl := fmt.Sprintf("http://localhost:%d/ping?logid=%s", app_const.SERVICE_PORT, logging.ValueLogID(c))

	ret, err := jaeger.JaegerSend(c, http.MethodGet, sendUrl, nil, nil, time.Second)
	ret, err = jaeger.JaegerSend(c, http.MethodGet, sendUrl, nil, nil, time.Second)
	if err != nil {
		fmt.Println(err)
	}

	sp, _ := jaeger.Inject(c, c.Request.Header, "get", jaeger.OPERATION_TYPE_REDIS)
	if sp != nil {
		defer sp.Finish()
	}

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
