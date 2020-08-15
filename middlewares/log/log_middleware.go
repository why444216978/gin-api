package log

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util/conversion"
	"gin-frame/libraries/util/dir"
	"gin-frame/libraries/util/random"
	"gin-frame/libraries/util/sys"
	"gin-frame/libraries/util/url"
	"gin-frame/libraries/xhop"

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

func LoggerMiddleware(port int, logFields map[string]string, productName, moduleName, env string) gin.HandlerFunc {
	runLogSection := "run"
	runLogConfig := config.GetConfig("log", runLogSection)
	runLogDir := runLogConfig.Key("dir").String()
	logArea, _ := runLogConfig.Key("area").Int()

	return func(c *gin.Context) {
		file := dir.CreateHourLogFile(runLogDir, moduleName+".log."+sys.HostName()+".")
		file = file + "/" + strconv.Itoa(random.RandomN(logArea))

		log.InitRun(&log.LogConfig{
			File:           file,
			Path:           runLogDir,
			Mode:           1,
			AsyncFormatter: false,
			Debug:          true,
		}, runLogDir, file)

		var logID string
		switch {
		case c.Query(logFields["query_id"]) != "":
			logID = c.Query(logFields["query_id"])
		case c.Request.Header.Get(logFields["header_id"]) != "":
			logID = c.Request.Header.Get(logFields["header_id"])
		default:
			logID = log.NewObjectId().Hex()
		}

		ctx := c.Request.Context()
		dst := new(log.LogFormat)

		dst.Port = port
		dst.LogId = logID
		dst.Method = c.Request.Method
		dst.CallerIp = c.ClientIP()
		dst.UriPath = c.Request.RequestURI
		dst.XHop = xhop.NextXhop(c.Request.Header, logFields["header_hop"])
		dst.Product = productName
		dst.Module = moduleName

		dst.Env = env

		ctx = log.ContextWithLogHeader(ctx, dst)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(logFields["header_id"], dst.LogId)
		c.Writer.Header().Set(logFields["header_hop"], dst.XHop.String())

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
