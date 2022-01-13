package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/library/logger"
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

type response struct {
	Code    uint64      `json:"code"`
	Toast   string      `json:"toast"`
	Data    interface{} `json:"data"`
	ErrMsg  string      `json:"errmsg"`
	TraceID string      `json:"trace_id"`
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
		Code:    code,
		Toast:   toast,
		Data:    data,
		ErrMsg:  errmsg,
		TraceID: logger.ValueTraceID(c.Request.Context()),
	})
	c.Abort()
}
