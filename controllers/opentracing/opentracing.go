package opentracing

import (
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

	ret, err := jaeger.JaegerSend(c, http.MethodPost, sendUrl, nil, nil, time.Second)
	if err != nil {
		fmt.Println(err)
	}

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Rpc1(c *gin.Context) {
	sendUrl := fmt.Sprintf("http://localhost:%d/test/conn?logid=%s", global.Global.AppPort, logging.ValueLogID(c))

	ret, err := jaeger.JaegerSend(c, http.MethodPost, sendUrl, nil, nil, time.Second)
	if err != nil {
		fmt.Println(err)
	}

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
