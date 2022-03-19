package api

import (
	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/app/response"
	gin_api "github.com/why444216978/gin-api/app/rpc/gin-api"
)

func Ping(c *gin.Context) {
	response.ResponseJSON(c, response.CodeSuccess, nil, "")
}

func RPC(c *gin.Context) {
	ret, err := gin_api.Ping(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, err.Error())
		return
	}
	response.ResponseJSON(c, response.CodeSuccess, ret, "")
}
