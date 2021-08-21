package log

import (
	"bytes"
	"encoding/json"
	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
	"go.uber.org/zap"
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

		ctx, span, traceID := jaeger.ExtractHTTP(ctx, c.Request, logging.ValueLogID(ctx))
		defer span.Finish()

		ctx = logging.WithTraceID(ctx, traceID)
		ctx = logging.AddTraceID(ctx, traceID)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		ctx = c.Request.Context()

		resp := responseWriter.body.String()
		respMap, _ := conversion.JsonToMap(resp)

		fields := logging.ValueHTTPFields(ctx)
		fields.Response = respMap
		fields.Code = c.Writer.Status()

		ctx = logging.WithHTTPRequestBody(ctx, fields.Request)

		req, _ := json.Marshal(fields.Request)
		jaeger.SetHTTPLog(span, string(req), resp)

		resource.ServiceLogger.Info("request info", zap.Reflect("data", fields))

		fields.Cost = int64(time.Now().Sub(start))

		c.Request = c.Request.WithContext(ctx)
	}
}
