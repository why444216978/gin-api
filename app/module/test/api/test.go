package api

import (
	"time"

	"github.com/why444216978/gin-api/app/response"
	gin_api "github.com/why444216978/gin-api/app/rpc/gin-api"
	"github.com/why444216978/go-util/http"

	"github.com/gin-gonic/gin"
)

func Rpc(c *gin.Context) {
	time.Sleep(time.Millisecond * 30)
	ret, err := gin_api.RPC(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, err.Error())
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, "")
}

type RPC1Request struct {
	A string `json:"a"`
}

func Rpc1(c *gin.Context) {
	time.Sleep(time.Millisecond * 99)
	var req RPC1Request
	if err := http.ParseAndValidateBody(c.Request, &req); err != nil {
		response.ResponseJSON(c, response.CodeParams, nil, err.Error())
		return
	}

	ret, err := gin_api.RPC1(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, err.Error())
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
