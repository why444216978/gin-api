package limiter

import (
	"gin-api/libraries/logging"
	"gin-api/resource"
	"gin-api/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
	"golang.org/x/time/rate"
)

func Limiter(maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second*1), maxBurstSize)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		fields := logging.ValueHTTPFields(c.Request.Context())
		fields.Response = map[string]interface{}{
			"code":   http.StatusServiceUnavailable,
			"toast":  "服务暂时不可用",
			"data":   "",
			"errmsg": "服务暂时不可用",
		}
		fields.Code = http.StatusInternalServerError

		data, _ := conversion.StructToMap(fields)
		resource.ServiceLogger.Error("panic", data) //这里不能打Fatal和Panic，否则程序会退出
		response.Response(c, response.CodeServer, nil, "")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
}
