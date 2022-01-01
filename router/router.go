package router

import (
	"net/http"
	"time"

	"github.com/why444216978/gin-api/api/conn"
	"github.com/why444216978/gin-api/api/ping"
	"github.com/why444216978/gin-api/api/test"
	"github.com/why444216978/gin-api/config"
	"github.com/why444216978/gin-api/library/middleware/limiter"
	"github.com/why444216978/gin-api/library/middleware/log"
	"github.com/why444216978/gin-api/library/middleware/panic"
	"github.com/why444216978/gin-api/library/middleware/timeout"
	"github.com/why444216978/gin-api/response"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	server := gin.New()

	server.Use(timeout.TimeoutMiddleware(time.Duration(config.App.ContextTimeout) * time.Millisecond))

	server.Use(panic.ThrowPanic())

	server.Use(log.InitContext())

	server.Use(limiter.Limiter(10))

	server.Use(log.LoggerMiddleware())

	server.NoRoute(func(c *gin.Context) {
		response.Response(c, response.CodeUriNotFound, nil, "")
		c.AbortWithStatus(http.StatusNotFound)
	})

	pingGroup := server.Group("/ping")
	{
		pingGroup.GET("", ping.Ping)
		pingGroup.GET("/rpc", ping.RPC)
	}

	testGroup := server.Group("/test")
	{
		testGroup.POST("/rpc", test.Rpc)
		testGroup.POST("/rpc1", test.Rpc1)
		testGroup.POST("/panic", test.Panic)
		testGroup.POST("/conn", conn.Do)
	}

	return server
}
