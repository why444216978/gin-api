package log

import "C"
import (
	"bytes"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"io/ioutil"
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

func WithContext() gin.HandlerFunc {
	logCfg := config.GetConfigToJson("log", "log")
	queryLogField := logCfg["query_field"].(string)
	headerLogField := logCfg["header_field"].(string)

	return func(c *gin.Context) {
		var logID string
		switch {
		case c.Query(queryLogField) != "":
			logID = c.Query(queryLogField)
		case c.Request.Header.Get(headerLogField) != "":
			logID = c.Request.Header.Get(headerLogField)
		default:
			logID = logging.NewObjectId().Hex()
		}
		c.Header(headerLogField, logID)

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		hostIP, _ := sys.ExternalIP()
		header := &logging.LogHeader{
			HTTPCode: c.Writer.Status(),
			Header:   c.Request.Header,
			LogId:    logID,
			CallerIp: c.ClientIP(),
			HostIp:   hostIP,
			Port:     app_const.SERVICE_PORT,
			UriPath:  c.Request.RequestURI,
			Module:   "http",
			Request:  reqBodyMap,
		}
		logging.WriteLogHeader(c, header)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next() // 处理请求

		responseBody := responseWriter.body.String()
		responseBodyMap, _ := conversion.JsonToMap(responseBody)

		header := logging.GetLogHeader(c)
		header.Cost = time.Now().Sub(start).Milliseconds()
		header.Response = responseBodyMap
		logging.WriteLogHeader(c, header)

		if !c.IsAborted() {
			logging.InfoCtx(c)
		}
	}
}
