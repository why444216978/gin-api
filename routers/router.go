package routers

import (
	"net/http"
	"time"

	"github.com/why444216978/gin-api/api/conn"
	"github.com/why444216978/gin-api/api/ping"
	"github.com/why444216978/gin-api/api/test"
	"github.com/why444216978/gin-api/config"
	"github.com/why444216978/gin-api/middlewares/limiter"
	"github.com/why444216978/gin-api/middlewares/log"
	"github.com/why444216978/gin-api/middlewares/panic"
	"github.com/why444216978/gin-api/middlewares/timeout"
	"github.com/why444216978/gin-api/response"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	server := gin.New()

	server.Use(log.InitContext())

	server.Use(panic.ThrowPanic())

	server.Use(timeout.TimeoutMiddleware(time.Duration(config.App.ContextTimeout) * time.Millisecond))

	server.Use(limiter.Limiter(10))

	server.Use(log.LoggerMiddleware())

	server.NoRoute(func(c *gin.Context) {
		response.Response(c, response.CodeUriNotFound, nil, "")
		c.AbortWithStatus(http.StatusNotFound)
	})

	pingGroup := server.Group("/ping")
	{
		pingGroup.GET("", ping.Ping)
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
