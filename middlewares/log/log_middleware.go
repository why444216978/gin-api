package log

import (
	"bytes"
	"gin-api/configs"
	"gin-api/libraries/config"
	"gin-api/libraries/util/random"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gin-api/libraries/log"
	"gin-api/libraries/util/conversion"
	"gin-api/libraries/util/dir"
	"gin-api/libraries/util/sys"
	"gin-api/libraries/util/url"
	"github.com/gin-gonic/gin"
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
	logDir, logArea := config.GetLogConfig(configs.LOG_SOURCE)
	return func(c *gin.Context) {
		file := dir.CreateHourLogFile(logDir, configs.SERVICE_NAME+".log."+sys.HostName()+".")
		file = file + "/" + strconv.Itoa(random.RandomN(logArea))

		log.InitRun(&log.LogConfig{
			File:           file,
			Path:           logDir,
			Mode:           1,
			AsyncFormatter: false,
			Debug:          true,
		}, logDir, file)

		var logID string
		switch {
		case c.Query(config.GetQueryLogIdField(configs.LOG_SOURCE)) != "":
			logID = c.Query(config.GetQueryLogIdField(configs.LOG_SOURCE))
		case c.Request.Header.Get(config.GetHeaderLogIdField(configs.LOG_SOURCE)) != "":
			logID = c.Request.Header.Get(config.GetHeaderLogIdField(configs.LOG_SOURCE))
		default:
			logID = log.NewObjectId().Hex()
		}

		ctx := c.Request.Context()
		dst := new(log.LogFormat)

		dst.Port = configs.SERVICE_PORT
		dst.LogId = logID
		dst.Method = c.Request.Method
		dst.CallerIp = c.ClientIP()
		dst.UriPath = c.Request.RequestURI
		dst.Product = configs.PRODUCT
		dst.Module = configs.MODULE
		dst.Env = configs.ENV

		ctx = log.ContextWithLogHeader(ctx, dst)
		c.Request = c.Request.WithContext(ctx)

		c.Header(config.GetHeaderLogIdField(configs.LOG_SOURCE), dst.LogId)

		c.Writer.Header().Set(config.GetHeaderLogIdField(configs.LOG_SOURCE), dst.LogId)

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		strReqBody := string(reqBody)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		dst.StartTime = time.Now()

		c.Next() // 处理请求

		dst.HttpCode = c.Writer.Status()

		responseBody := responseWriter.body.String()

		if dst.HttpCode == http.StatusOK {
			log.Info(dst, map[string]interface{}{
				"requestHeader": c.Request.Header,
				"requestBody":   conversion.JsonToMap(strReqBody),
				"responseBody":  conversion.JsonToMap(responseBody),
				"uriQuery":      url.ParseUriQueryToMap(c.Request.URL.RawQuery),
			})
		}
	}
}
