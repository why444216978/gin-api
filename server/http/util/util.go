package util

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

//定义新的struck，继承gin的ResponseWriter
//添加body字段，用于将response暴露给日志
type BodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

//gin的ResponseWriter继承的底层http server
//实现http的Write方法，额外添加一个body字段，用于获取response body
func (w BodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
