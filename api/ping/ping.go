package ping

import (
	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/response"
	gin_api "github.com/why444216978/gin-api/rpc/gin-api"
)

func Ping(c *gin.Context) {
	response.Response(c, response.CodeSuccess, nil, "")
}
func RPC(c *gin.Context) {
	ret, err := gin_api.Ping(c.Request.Context())
	if err != nil {
		response.Response(c, response.CodeServer, ret, err.Error())
		return
	}
	response.Response(c, response.CodeSuccess, nil, "")
}
