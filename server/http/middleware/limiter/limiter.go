package limiter

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/server/http/response"
)

func Limiter(maxBurstSize int, l logger.Logger) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second*1), maxBurstSize)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		fields := logger.ValueFields(c.Request.Context())
		fields.Response = map[string]interface{}{
			"code":   http.StatusServiceUnavailable,
			"toast":  "服务暂时不可用",
			"data":   "",
			"errmsg": "服务暂时不可用",
		}
		fields.Code = http.StatusInternalServerError
		ctx = logger.WithFields(ctx, fields)

		c.Request = c.Request.WithContext(ctx)

		l.Error(ctx, "panic") // 这里不能打Fatal和Panic，否则程序会退出
		response.ResponseJSON(c, http.StatusServiceUnavailable, nil, response.WrapToast(nil, http.StatusText(http.StatusServiceUnavailable)))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
