package panic

import (
	"bytes"
	"gin-api/codes"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/util/conversion"
	"gin-api/libraries/util/sys"
	"gin-api/libraries/util/url"
	"gin-api/libraries/logging"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
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
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
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

				var logId string
				switch {
				case c.Query(config.GetQueryLogIdField(app_const.LOG_SOURCE)) != "":
					logId = c.Query(config.GetQueryLogIdField(app_const.LOG_SOURCE))
				case c.Request.Header.Get(config.GetHeaderLogIdField(app_const.LOG_SOURCE)) != "":
					logId = c.Request.Header.Get(config.GetHeaderLogIdField(app_const.LOG_SOURCE))
				default:
					logId = logging.NewObjectId().Hex()
				}

				c.Header(config.GetHeaderLogIdField(app_const.LOG_SOURCE), logId)

				c.Writer.Header().Set(config.GetHeaderLogIdField(app_const.LOG_SOURCE), logId)

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

				hostIp,_ := sys.ExternalIP()

				header := &logging.LogHeader{
					LogId: logId,
					CallerIp: c.ClientIP(),
					HostIp: hostIp,
					Port: app_const.SERVICE_PORT,
					Product: app_const.PRODUCT,
					Module: app_const.MODULE,
					ServiceId: app_const.SERVICE_NAME,
					UriPath: c.Request.RequestURI,
					Env: app_const.ENV,
				}

				logging.Error(header, map[string]interface{}{
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
