package limiter

import (
	"net/http"
	"time"

	"github.com/why444216978/gin-api/library/logging"
	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Limiter(maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second*1), maxBurstSize)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		fields := logging.ValueHTTPFields(c.Request.Context())
		fields.Response = map[string]interface{}{
			"code":   http.StatusServiceUnavailable,
			"toast":  "服务暂时不可用",
			"data":   "",
			"errmsg": "服务暂时不可用",
		}
		fields.Code = http.StatusInternalServerError
		ctx = logging.WithHTTPFields(ctx, fields)

		c.Request = c.Request.WithContext(ctx)

		resource.ServiceLogger.Error(ctx, "panic") //这里不能打Fatal和Panic，否则程序会退出
		response.Response(c, response.CodeUnavailable, nil, "")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
}
