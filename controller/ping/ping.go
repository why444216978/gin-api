package ping

import (
	"github.com/why444216978/gin-api/response"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	response.Response(c, response.CodeSuccess, nil, "")
}
