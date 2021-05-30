package panic

import (
	"bytes"
	"gin-api/app_const"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"gin-api/response"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"
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

				common := &logging.Common{
					LogID: logging.GetLogID(c),
				}
				logging.WriteLogCommon(c, common)

				hostIP, _ := sys.ExternalIP()

				// `{"code":500,"toast":"服务器错误","data":{},"errmsg":""}`
				fields := logging.Fields{
					Header:  c.Request.Header,
					Method:  c.Request.Method,
					Request: logging.GetRequestBody(c),
					Response: map[string]interface{}{
						"code":   http.StatusInternalServerError,
						"toast":  "服务器错误",
						"data":   "",
						"errmsg": "",
					},
					Code:     http.StatusInternalServerError,
					CallerIP: c.ClientIP(),
					HostIP:   hostIP,
					Port:     app_const.SERVICE_PORT,
					API:      c.Request.RequestURI,
					Module:   "HTTP",
					Trace:    debugStack,
				}
				fields.Common = *common

				data, _ := conversion.StructToMap(fields)
				resource.Logger.Error("panic", data) //这里不能打Fatal和Panic，否则程序会退出
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
