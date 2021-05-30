package log

import "C"
import (
	"bytes"
	"gin-api/app_const"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"
)

//定义新的struck，继承gin的ResponseWriter
//添加body字段，用于将response暴露给日志
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

//gin的ResponseWriter继承的底层http server
//实现http的Write方法，额外添加一个body字段，用于获取response body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

const (
	LOG_HEADER = "X-Log-Id"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next()

		resp := responseWriter.body.String()
		respMap, _ := conversion.JsonToMap(resp)

		common := &logging.Common{
			LogID: logging.GetLogID(c),
		}
		logging.WriteLogCommon(c, common)

		hostIP, _ := sys.ExternalIP()

		fields := logging.Fields{
			Header:   c.Request.Header,
			Method:   c.Request.Method,
			Request:  logging.GetRequestBody(c),
			Response: respMap,
			Code:     c.Writer.Status(),
			CallerIP: c.ClientIP(),
			HostIP:   hostIP,
			Port:     app_const.SERVICE_PORT,
			API:      c.Request.RequestURI,
			Module:   "HTTP",
			Cost:     int64(time.Now().Sub(start)),
		}
		fields.Common = *common

		data, _ := conversion.StructToMap(fields)
		resource.Logger.Info("request info", data)
	}
}
