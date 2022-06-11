package api

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/app/response"
	gin_api "github.com/why444216978/gin-api/app/rpc/gin-api"
	httpResponse "github.com/why444216978/gin-api/server/http/response"
	"github.com/why444216978/go-util/http"
)

func Rpc(c *gin.Context) {
	time.Sleep(time.Millisecond * 30)
	ret, err := gin_api.RPC(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}

type RPC1Request struct {
	A string `json:"a"`
}

func Rpc1(c *gin.Context) {
	time.Sleep(time.Millisecond * 99)
	var req RPC1Request
	if err := http.ParseAndValidateBody(c.Request, &req); err != nil {
		response.ResponseJSON(c, response.CodeParams, nil, httpResponse.WrapToast(err, err.Error()))
		return
	}

	ret, err := gin_api.RPC1(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}

func Panic(c *gin.Context) {
	panic("test err")
}
