package timeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware  超时控制
func TimeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
