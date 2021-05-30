package limiter

import (
	"fmt"
	"net/http"
	"time"

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
		fmt.Println("Too many requests")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"errno":    http.StatusServiceUnavailable,
			"errmsg":   "服务暂时不可用",
			"data":     nil,
			"user_msg": "服务暂时不可用",
		})
		c.Abort()
		return
	}
}
