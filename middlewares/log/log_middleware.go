package log

import "C"
import (
	"bytes"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"
	"github.com/why444216978/go-util/url"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

func LoggerMiddleware() gin.HandlerFunc {
	envCfg := config.GetConfigToJson("env", "env")
	logCfg := config.GetConfigToJson("log", "log")
	queryLogField := logCfg["query_field"].(string)
	headerLogField := logCfg["header_field"].(string)

	return func(c *gin.Context) {
		var logId string
		switch {
		case c.Query(queryLogField) != "":
			logId = c.Query(queryLogField)
		case c.Request.Header.Get(headerLogField) != "":
			logId = c.Request.Header.Get(headerLogField)
		default:
			logId = logging.NewObjectId().Hex()
		}

		c.Header(headerLogField, logId)

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		strReqBody := string(reqBody)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next() // 处理请求

		responseBody := responseWriter.body.String()

		hostIp, _ := sys.ExternalIP()

		header := &logging.LogHeader{
			LogId:     logId,
			CallerIp:  c.ClientIP(),
			HostIp:    hostIp,
			Port:      app_const.SERVICE_PORT,
			Product:   app_const.PRODUCT,
			Module:    app_const.MODULE,
			ServiceId: app_const.SERVICE_NAME,
			UriPath:   c.Request.RequestURI,
			Env:       envCfg["env"].(string),
		}
		logging.Info(header, map[string]interface{}{
			"requestHeader": c.Request.Header,
			"requestBody":   conversion.JsonToMap(strReqBody),
			"responseBody":  conversion.JsonToMap(responseBody),
			"uriQuery":      url.ParseUriQueryToMap(c.Request.URL.RawQuery),
			"http_code":     c.Writer.Status(),
		})
	}
}
