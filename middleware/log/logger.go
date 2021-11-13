package log

import (
	"bytes"
	"time"

	jaeger_http "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logging"
	"github.com/why444216978/gin-api/resource"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
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
	return func(c *gin.Context) {
		start := time.Now()

		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		ctx := c.Request.Context()

		ctx, span, traceID := jaeger_http.ExtractHTTP(ctx, c.Request, logging.ValueLogID(ctx))
		defer span.Finish()

		ctx = logging.WithTraceID(ctx, traceID)
		ctx = logging.AddTraceID(ctx, traceID)

		//这里需要写入ctx，否则会断开
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		fields := logging.ValueHTTPFields(ctx)

		//resp处理
		resp := responseWriter.body.String()
		respMap, _ := conversion.JsonToMap(resp)

		//span写入req和resp
		req, _ := conversion.JsonEncode(fields.Request)
		jaeger_http.SetHTTPLog(span, string(req), resp)

		//追加fields
		fields.Response = respMap
		fields.Code = c.Writer.Status()
		fields.Cost = time.Since(start).Milliseconds()

		ctx = logging.WithHTTPFields(ctx, fields)

		resource.ServiceLogger.Info(ctx, "request info")

		c.Request = c.Request.WithContext(ctx)
	}
}
