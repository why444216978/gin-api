package panic

import (
	"bytes"
	"gin-api/libraries/logging"
	"gin-api/response"
	"net/http"
	"runtime/debug"
	"strings"

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
	// logCfg := config.GetConfigToJson("log", "log")
	// queryLogField := logCfg["query_field"].(string)
	// headerLogField := logCfg["header_field"].(string)
	return func(c *gin.Context) {
		// var logID string
		// switch {
		// case c.Query(queryLogField) != "":
		// 	logID = c.Query(queryLogField)
		// case c.Request.Header.Get(headerLogField) != "":
		// 	logID = c.Request.Header.Get(headerLogField)
		// default:
		// 	logID = logging.NewObjectId().Hex()
		// }
		// c.Header(headerLogField, logID)

		// reqBody := []byte{}
		// if c.Request.Body != nil { // Read
		// 	reqBody, _ = ioutil.ReadAll(c.Request.Body)
		// }
		// reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
		// c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		// hostIP, _ := sys.ExternalIP()
		// header := &logging.LogHeader{
		// 	HTTPCode: c.Writer.Status(),
		// 	Header:   c.Request.Header,
		// 	LogId:    logID,
		// 	CallerIp: c.ClientIP(),
		// 	HostIp:   hostIP,
		// 	Port:     app_const.SERVICE_PORT,
		// 	UriPath:  c.Request.RequestURI,
		// 	Module:   "http",
		// 	Request:  reqBodyMap,
		// }
		// logging.WriteLogHeader(c, header)

		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				mailDebugStack := ""
				debugStack := make(map[int]interface{})
				for k, v := range strings.Split(string(debug.Stack()), "\n") {
					//fmt.Println(v)
					mailDebugStack += v + "<br>"
					debugStack[k] = v
				}

				responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
				c.Writer = responseWriter

				header := logging.GetLogHeader(c)
				header.HTTPCode = http.StatusInternalServerError
				header.Trace = debugStack
				header.Error = err
				logging.WriteLogHeader(c, header)

				logging.ErrorCtx(c)
				response.Response(c, response.CODE_SERVER, nil, "")
				c.AbortWithStatus(http.StatusInternalServerError)

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
			}
		}(c)
		c.Next()
	}
}
