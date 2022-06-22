package log

import (
	"bytes"
	"encoding/base64"
	"net/http/httputil"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/assert"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"

	"github.com/why444216978/gin-api/library/app"
	jaegerHTTP "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/server/http/util"
)

func LoggerMiddleware(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		defer func() {
			c.Request = c.Request.WithContext(ctx)
		}()

		start := time.Now()

		serverIP, _ := sys.LocalIP()

		logID := logger.ExtractLogID(c.Request)
		ctx = logger.WithLogID(ctx, logID)

		// req := logger.GetRequestBody(c.Request)
		req, _ := httputil.DumpRequest(c.Request, true)

		responseWriter := &util.BodyWriter{Body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		ctx, span, traceID := jaegerHTTP.ExtractHTTP(ctx, c.Request, logger.ValueLogID(ctx))
		if !assert.IsNil(span) {
			defer span.Finish()
		}
		ctx = logger.WithTraceID(ctx, traceID)

		fields := []logger.Field{
			logger.Reflect(logger.LogID, logID),
			logger.Reflect(logger.TraceID, traceID),
			logger.Reflect(logger.Header, c.Request.Header),
			logger.Reflect(logger.Method, c.Request.Method),
			logger.Reflect(logger.Request, base64.StdEncoding.EncodeToString(req)),
			logger.Reflect(logger.Response, make(map[string]interface{})),
			logger.Reflect(logger.ClientIP, c.ClientIP()),
			logger.Reflect(logger.ClientPort, 0),
			logger.Reflect(logger.ServerIP, serverIP),
			logger.Reflect(logger.ServerPort, app.Port()),
			logger.Reflect(logger.API, c.Request.RequestURI),
		}
		// Next之前这里需要写入ctx，否则会丢失log、断开trace
		ctx = logger.WithFields(ctx, fields)
		c.Request = c.Request.WithContext(ctx)

		var doneFlag int32
		done := make(chan struct{}, 1)
		defer func() {
			done <- struct{}{}
			atomic.StoreInt32(&doneFlag, 1)

			resp := responseWriter.Body.Bytes()
			respString := string(resp)
			if responseWriter.Body.Len() > 0 {
				logResponse, _ := conversion.JsonToMap(respString)
				ctx = logger.AddField(ctx, logger.Reflect(logger.Response, logResponse))
			}

			reqString, _ := conversion.JsonEncode(req)
			jaegerHTTP.SetHTTPLog(span, reqString, respString)

			ctx = logger.AddField(ctx,
				logger.Reflect(logger.Code, c.Writer.Status()),
				logger.Reflect(logger.Cost, time.Since(start).Milliseconds()),
			)
			l.Info(ctx, "request info")
		}()

		go func() {
			select {
			case <-done:
			case <-ctx.Done():
				if atomic.LoadInt32(&doneFlag) == 1 {
					return
				}
				ctx = logger.AddField(ctx,
					logger.Reflect(logger.Code, 499),
					logger.Reflect(logger.Cost, time.Since(start).Milliseconds()),
				)
				l.Warn(ctx, "client canceled")
			}
		}()

		c.Next()
	}
}
