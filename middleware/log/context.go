package log

import (
	app_config "github.com/why444216978/gin-api/config"
	"github.com/why444216978/gin-api/library/logging"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/sys"
)

func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverIP, _ := sys.LocalIP()

		logID := logging.ExtractLogID(c.Request)

		req := logging.GetRequestBody(c.Request)
		fields := logging.Fields{
			LogID:      logID,
			Header:     c.Request.Header,
			Method:     c.Request.Method,
			Request:    req,
			ClientIP:   c.ClientIP(),
			ClientPort: 0,
			ServerIP:   serverIP,
			ServerPort: app_config.App.AppPort,
			API:        c.Request.RequestURI,
			Module:     logging.ModuleHTTP,
		}

		ctx := logging.WithLogID(c.Request.Context(), logID)
		ctx = logging.WithHTTPRequestBody(ctx, req)
		ctx = logging.WithHTTPFields(ctx, fields)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
