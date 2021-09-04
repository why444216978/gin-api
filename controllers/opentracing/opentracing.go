package opentracing

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	lib_http "gin-api/libraries/http"
	"gin-api/libraries/logging"
	"gin-api/response"

	"github.com/gin-gonic/gin"
)

func Rpc(c *gin.Context) {
	uri := fmt.Sprintf("/test/rpc1?logid=%s", logging.ValueLogID(c))

	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := lib_http.Send(c.Request.Context(), "gin-api", http.MethodPost, uri, header, bytes.NewBufferString(`{"a":"a"}`), time.Second)
	if err != nil {
		fmt.Println(ret)
		fmt.Println(err)
		return
	}

	response.Response(c, response.CodeSuccess, ret, "")
}

func Rpc1(c *gin.Context) {
	uri := fmt.Sprintf("/test/conn?logid=%s", logging.ValueLogID(c))

	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := lib_http.Send(c.Request.Context(), "gin-api", http.MethodPost, uri, header, nil, time.Second)
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
