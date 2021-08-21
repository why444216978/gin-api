package panic

import (
	"bytes"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"gin-api/response"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
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
				mailDebugStack := ""
				debugStack := make(map[int]interface{})
				for k, v := range strings.Split(string(debug.Stack()), "\n") {
					//fmt.Println(v)
					mailDebugStack += v + "<br>"
					debugStack[k] = v
				}

				fields := logging.ValueHTTPFields(c.Request.Context())
				fields.Response = map[string]interface{}{
					"code":   http.StatusInternalServerError,
					"toast":  "服务器错误",
					"data":   "",
					"errmsg": "服务器错误",
				}
				fields.Code = http.StatusInternalServerError
				fields.Trace = debugStack

				data, _ := conversion.StructToMap(fields)
				resource.ServiceLogger.Error("panic", data) //这里不能打Fatal和Panic，否则程序会退出
				response.Response(c, response.CodeServer, nil, "")
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
