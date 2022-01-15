package log

import (
	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/sys"

	appConfig "github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/library/logger"
	loggerHTTP "github.com/why444216978/gin-api/library/logger/http"
)

func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverIP, _ := sys.LocalIP()

		logID := loggerHTTP.ExtractLogID(c.Request)

		req := loggerHTTP.GetRequestBody(c.Request)
		fields := logger.Fields{
			LogID:      logID,
			Header:     c.Request.Header,
			Method:     c.Request.Method,
			Request:    req,
			ClientIP:   c.ClientIP(),
			ClientPort: 0,
			ServerIP:   serverIP,
			ServerPort: appConfig.App.AppPort,
			API:        c.Request.RequestURI,
		}

		ctx := logger.WithLogID(c.Request.Context(), logID)
		ctx = logger.WithHTTPRequestBody(ctx, req)
		ctx = logger.WithHTTPFields(ctx, fields)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
