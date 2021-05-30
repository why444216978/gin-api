package timeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware  超时控制
func TimeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// defer func() {
		// 	// check if context timeout was reached
		// 	if ctx.Err() == context.DeadlineExceeded {
		// 		// write response and abort the request
		// 		// c.Writer.WriteHeader(http.StatusGatewayTimeout)
		// 		response.Response(c, response.CODE_SERVER, nil, "")
		// 		c.Abort()
		// 	}
		// 	//cancel to clear resources after finished
		// 	cancel()
		// }()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
