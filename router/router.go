package router

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/config"
	"github.com/why444216978/gin-api/library/middleware/limiter"
	"github.com/why444216978/gin-api/library/middleware/log"
	"github.com/why444216978/gin-api/library/middleware/panic"
	"github.com/why444216978/gin-api/library/middleware/timeout"
	conn "github.com/why444216978/gin-api/module/goods/api"
	ping "github.com/why444216978/gin-api/module/ping/api"
	test "github.com/why444216978/gin-api/module/test/api"
	"github.com/why444216978/gin-api/response"
)

func InitRouter() *gin.Engine {
	server := gin.New()

	startPprof(server)

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

func startPprof(server *gin.Engine) {
	if !config.App.Pprof {
		return
	}
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	pprof.Register(server)
}
