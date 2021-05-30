package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Code uint64

const (
	CODE_SUCCESS Code = 0
	CODE_PARAMS  Code = 1
	CODE_SERVER  Code = 2
	CODE_TIMEOUT Code = 3
	CODE_URI     Code = 4
)

var codeText = map[Code]string{
	CODE_SUCCESS: "success",
	CODE_PARAMS:  "参数错误",
	CODE_SERVER:  "服务器错误",
	CODE_TIMEOUT: "超时服务端关闭",
	CODE_URI:     "资源不存在",
}

type response struct {
	Code   Code        `json:"code"`
	Toast  string      `json:"toast"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
}

func Response(c *gin.Context, code Code, data interface{}, errmsg string) {
	if data == nil {
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
	c.Done()
}
