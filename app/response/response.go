package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/server/http/response"
)

const (
	CodeSuccess     = 0
	CodeParams      = 1
	CodeUriNotFound = http.StatusNotFound
	CodeServer      = http.StatusInternalServerError
	CodeUnavailable = http.StatusServiceUnavailable
	CodeTimeout     = http.StatusGatewayTimeout
)

var codeText = map[uint64]string{
	CodeSuccess:     "success",
	CodeParams:      "参数错误",
	CodeUriNotFound: "资源不存在",
	CodeUnavailable: "服务器暂时不可用",
	CodeTimeout:     "请求超时",
	CodeServer:      "服务器错误",
}

func ResponseJSON(c *gin.Context, code uint64, data interface{}, errmsg string) {
	toast, ok := codeText[code]
	if !ok {
		toast = ""
	}

	response.ResponseJSON(c, code, data, errmsg, toast)
}
