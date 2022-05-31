package api

import (
	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/app/response"
	gin_api "github.com/why444216978/gin-api/app/rpc/gin-api"
	httpResponse "github.com/why444216978/gin-api/server/http/response"
)

func Ping(c *gin.Context) {
	response.ResponseJSON(c, response.CodeSuccess, nil, nil)
}

func RPC(c *gin.Context) {
	ret, err := gin_api.Ping(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}
	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}
