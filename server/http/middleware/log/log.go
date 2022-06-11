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

	appConfig "github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/resource"
	jaegerHTTP "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logger"
	loggerHTTP "github.com/why444216978/gin-api/library/logger/http"
	"github.com/why444216978/gin-api/server/http/util"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		defer func() {
			c.Request = c.Request.WithContext(ctx)
		}()

		start := time.Now()

		serverIP, _ := sys.LocalIP()

		logID := loggerHTTP.ExtractLogID(c.Request)
		ctx = logger.WithLogID(ctx, logID)

		// req := loggerHTTP.GetRequestBody(c.Request)
		req, _ := httputil.DumpRequest(c.Request, true)

		responseWriter := &util.BodyWriter{Body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		ctx, span, traceID := jaegerHTTP.ExtractHTTP(ctx, c.Request, logger.ValueLogID(ctx))
		if !assert.IsNil(span) {
			defer span.Finish()
		}
		ctx = logger.WithTraceID(ctx, traceID)

		fields := logger.Fields{
			LogID:      logID,
			TraceID:    traceID,
			Header:     c.Request.Header,
			Method:     c.Request.Method,
			Request:    base64.StdEncoding.EncodeToString(req),
			Response:   make(map[string]interface{}),
			ClientIP:   c.ClientIP(),
			ClientPort: 0,
			ServerIP:   serverIP,
			ServerPort: appConfig.App.AppPort,
			API:        c.Request.RequestURI,
		}
		// Next之前这里需要写入ctx，否则会丢失log、断开trace
		ctx = logger.WithHTTPFields(ctx, fields)
		c.Request = c.Request.WithContext(ctx)

		var doneFlag int32
		done := make(chan struct{}, 1)
		defer func() {
			done <- struct{}{}
			atomic.StoreInt32(&doneFlag, 1)

			resp := responseWriter.Body.Bytes()
			respString := string(resp)
			if responseWriter.Body.Len() > 0 {
				fields.Response, _ = conversion.JsonToMap(respString)
			}

			reqString, _ := conversion.JsonEncode(req)
			jaegerHTTP.SetHTTPLog(span, reqString, respString)

			fields.Code = c.Writer.Status()
			fields.Cost = time.Since(start).Milliseconds()
			ctx = logger.WithHTTPFields(ctx, fields)
			resource.ServiceLogger.Info(ctx, "request info")
		}()

		go func() {
			select {
			case <-done:
			case <-ctx.Done():
				if atomic.LoadInt32(&doneFlag) == 1 {
					return
				}
				fields.Code = 499
				fields.Cost = time.Since(start).Milliseconds()
				ctx = logger.WithHTTPFields(ctx, fields)
				resource.ServiceLogger.Warn(ctx, "client canceled")
			}
		}()

		c.Next()
	}
}
