package opentracing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/why444216978/gin-api/libraries/logging"
	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/response"
	"github.com/why444216978/go-util/validate"

	"github.com/gin-gonic/gin"
)

func Rpc(c *gin.Context) {
	uri := fmt.Sprintf("/test/rpc1?logid=%s", logging.ValueLogID(c))

	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := resource.HTTPRPC.Send(c.Request.Context(), "gin-api", http.MethodPost, uri, header, bytes.NewBufferString(`{"a":"a"}`), time.Second)
	if err != nil {
		fmt.Println(ret)
		fmt.Println(err)
		return
	}

	response.Response(c, response.CodeSuccess, ret, "")
}

type RPC1Request struct {
	A string `json:"a"`
}

func Rpc1(c *gin.Context) {
	var req RPC1Request
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.Response(c, response.CodeParams, nil, err.Error())
		return
	}
	if err := validate.Validate(&req); err != nil {
		response.Response(c, response.CodeParams, nil, err.Error())
		return
	}
	fmt.Println(req)

	uri := fmt.Sprintf("/test/conn?logid=%s", logging.ValueLogID(c))
	header := map[string]string{logging.LogHeader: logging.ValueLogID(c.Request.Context())}
	ret, err := resource.HTTPRPC.Send(c.Request.Context(), "gin-api", http.MethodPost, uri, header, nil, time.Second)
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
