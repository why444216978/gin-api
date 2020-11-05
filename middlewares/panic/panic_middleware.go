package panic

import (
	"bytes"
	"fmt"
	"gin-api/codes"
	"gin-api/configs"
	"gin-api/libraries/config"
	"gin-api/libraries/util/dir"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"gin-api/libraries/log"
	"gin-api/libraries/util"
	"gin-api/libraries/util/conversion"
	"gin-api/libraries/util/url"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ThrowPanic() gin.HandlerFunc {
	logDir, logArea := config.GetLogConfig(configs.LOG_SOURCE)
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"errno":    codes.SERVER_ERROR,
					"errmsg":   codes.ErrorMsg[codes.SERVER_ERROR],
					"data":     make(map[string]interface{}),
					"user_msg": codes.ErrorUserMsg[codes.SERVER_ERROR],
				})

				mailDebugStack := ""
				debugStack := make(map[int]interface{})
				for k, v := range strings.Split(string(debug.Stack()), "\n") {
					//fmt.Println(v)
					mailDebugStack += v + "<br>"
					debugStack[k] = v
				}

				file := dir.CreateHourLogFile(logDir, configs.SERVICE_NAME+".err."+util.HostName()+".")
				file = file + "/" + strconv.Itoa(util.RandomN(logArea))

				log.InitError(&log.LogConfig{
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

				logHeader := &log.LogFormat{}
				ctx := c.Request.Context()
				dst := new(log.LogFormat)
				*dst = *logHeader

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

				dst.HttpCode = c.Writer.Status()

				responseBody := responseWriter.body.String()

				log.Error(dst, map[string]interface{}{
					"requestHeader": c.Request.Header,
					"requestBody":   conversion.JsonToMap(strReqBody),
					"responseBody":  conversion.JsonToMap(responseBody),
					"uriQuery":      url.ParseUriQueryToMap(c.Request.URL.RawQuery),
					"err":           err,
					"trace":         debugStack,
				})

				/* util.WriteWithIo(file,"[" +dateTime+"]")
				util.WriteWithIo(file, fmt.Sprintf("%v\r\n", err))
				util.WriteWithIo(file, debugStack) */

				//subject := fmt.Sprintf("【重要错误】%s 项目出错了！", "go-gin")
				//
				//body := strings.ReplaceAll(MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", err))
				//body = strings.ReplaceAll(body, "{RequestTime}", util_time.GetCurrentDate())
				//body = strings.ReplaceAll(body, "{RequestURL}", c.Request.Method+"  "+c.Request.Host+c.Request.RequestURI)
				//body = strings.ReplaceAll(body, "{RequestUA}", c.Request.UserAgent())
				//body = strings.ReplaceAll(body, "{RequestIP}", c.ClientIP())
				//body = strings.ReplaceAll(body, "{DebugStack}", mailDebugStack)
				//
				//options := &mail.Options{
				//	MailHost: "smtp.163.com",
				//	MailPort: 465,
				//	MailUser: "weihaoyu@163.com",
				//	MailPass: "",
				//	MailTo:   "weihaoyu@163.com",
				//	Subject:  subject,
				//	Body:     body,
				//}
				//_ = mail.Send(options)

				c.Done()
			}
		}(c)
		c.Next()
	}
}
