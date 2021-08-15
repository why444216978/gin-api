package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess     = 0
	CodeParams      = 1
	CodeUriNotFound = http.StatusNotFound
	CodeServer      = http.StatusInternalServerError
	CodeTimeout     = http.StatusGatewayTimeout
)

var codeText = map[uint64]string{
	CodeSuccess:     "success",
	CodeParams:      "参数错误",
	CodeUriNotFound: "资源不存在",
	CodeServer:      "服务器错误",
	CodeTimeout:     "请求超时",
}

type response struct {
	Code   uint64      `json:"code"`
	Toast  string      `json:"toast"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
}

func Response(c *gin.Context, code uint64, data interface{}, errmsg string) {
	if data == nil || code != CodeSuccess {
		data = make(map[string]interface{})
	}

	toast, ok := codeText[code]
	if !ok {
		toast = ""
	}
	c.JSON(http.StatusOK, response{
		Code:   code,
		Toast:  toast,
		Data:   data,
		ErrMsg: errmsg,
	})
	c.Abort()
}
