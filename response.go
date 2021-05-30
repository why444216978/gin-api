package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Code uint64

const (
	CODE_SUCCESS Code = 0
	CODE_PARAMS  Code = 1
	CODE_SERVER  Code = 2
)

var codeText = map[Code]string{
	CODE_SUCCESS: "success",
	CODE_PARAMS:  "参数错误",
	CODE_SERVER:  "服务器错误",
}

type response struct {
	Code   Code        `json:"code"`
	Toast  string      `json:"toast"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
}

func Response(c *gin.Context, code Code, data interface{}, errmsg string) {
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
