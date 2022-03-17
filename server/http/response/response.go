package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/assert"

	"github.com/why444216978/gin-api/library/logger"
)

type response struct {
	Code    uint64      `json:"code"`
	Toast   string      `json:"toast"`
	Data    interface{} `json:"data"`
	ErrMsg  string      `json:"errmsg"`
	TraceID string      `json:"trace_id"`
}

func ResponseJSON(c *gin.Context, code uint64, data interface{}, errmsg, toast string) {
	if assert.IsNil(data) {
		data = make(map[string]interface{})
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
