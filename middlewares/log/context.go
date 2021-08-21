package log

import (
	"gin-api/libraries/logging"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/sys"
)

func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		hostIP, _ := sys.ExternalIP()

		logID := logging.ExtractLogID(c.Request)

		fields := logging.Fields{
			LogID:    logID,
			Header:   c.Request.Header,
			Method:   c.Request.Method,
			Request:  logging.GetRequestBody(c.Request),
			CallerIP: c.ClientIP(),
			HostIP:   hostIP,
			API:      c.Request.RequestURI,
			Module:   logging.ModuleHTTP,
		}

		ctx := logging.WithLogID(c.Request.Context(), logID)
		ctx = logging.WithHTTPRequestBody(ctx, fields.Request)
		ctx = logging.WithHTTPFields(ctx, fields)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
