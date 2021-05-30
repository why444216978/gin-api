package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CODE_SUCCESS       = 0
	CODE_PARAMS        = 1
	CODE_URI_NOT_FOUND = http.StatusNotFound
	CODE_SERVER        = http.StatusInternalServerError
	CODE_TIMEOUT       = http.StatusGatewayTimeout
)

var codeText = map[uint64]string{
	CODE_SUCCESS:       "success",
	CODE_PARAMS:        "参数错误",
	CODE_URI_NOT_FOUND: "资源不存在",
	CODE_SERVER:        "服务器错误",
	CODE_TIMEOUT:       "请求超时",
}

type response struct {
	Code   uint64      `json:"code"`
	Toast  string      `json:"toast"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
}

func Response(c *gin.Context, code uint64, data interface{}, errmsg string) {
	if data == nil || code != CODE_SUCCESS {
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
	c.Next()
}
