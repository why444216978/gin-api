package timeout

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TimeoutKey = "Timeout-Millisecond"
	startKey   = "Timeout-StartAt"
)

// TimeoutMiddleware  超时控制中间件
func TimeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		remain := timeout
		headerTimeout := c.Request.Header.Get(TimeoutKey)
		if headerTimeout != "" {
			t, _ := strconv.ParseInt(headerTimeout, 10, 64)
			remain = time.Duration(t) * time.Millisecond
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), remain)
		_ = cancel

		ctx = SetStart(ctx, remain.Milliseconds())

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func SetStart(ctx context.Context, timeout int64) context.Context {
	ctx = context.WithValue(ctx, TimeoutKey, timeout)
	return context.WithValue(ctx, startKey, nowMillisecond())
}

func CalcRemainTimeout(ctx context.Context) (int64, error) {
	timeout, ok := ctx.Value(TimeoutKey).(int64)
	if !ok {
		return 0, nil
	}

	startAt, ok := ctx.Value(startKey).(int64)
	if !ok {
		return 0, errors.New("miss startAt")
	}

	remain := timeout - (nowMillisecond() - startAt)
	if remain < 0 {
		return 0, errors.New("timeout < diff, context deadline exceeded")
	}

	return remain, nil
}

func nowMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
