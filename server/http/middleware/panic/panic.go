package panic

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/server/http/response"
)

func ThrowPanic(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				// mailDebugStack := ""
				// debugStack := make(map[int]interface{})
				// for k, v := range strings.Split(string(debug.Stack()), "\n") {
				// 	//fmt.Println(v)
				// 	mailDebugStack += v + "<br>"
				// 	debugStack[k] = v
				// }

				ctx := c.Request.Context()

				fields := logger.ValueFields(ctx)
				ctx = logger.AddField(ctx,
					logger.Reflect(logger.Code, http.StatusInternalServerError),
					logger.Reflect(logger.Response, map[string]interface{}{
						"code":   http.StatusInternalServerError,
						"toast":  "服务器错误",
						"data":   "",
						"errmsg": "服务器错误",
					}),
					// logger.Reflect(logger.Trace, debugStack),
				)
				ctx = logger.WithFields(c.Request.Context(), fields)
				c.Request = c.Request.WithContext(ctx)

				l.Error(ctx, fmt.Sprintf("%s", err)) // 这里不能打Fatal和Panic，否则程序会退出
				response.ResponseJSON(c, http.StatusInternalServerError, nil, response.WrapToast(nil, http.StatusText(http.StatusInternalServerError)))
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
