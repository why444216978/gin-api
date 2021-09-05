package routers

import (
	"net/http"
	"time"

	"gin-api/config"
	"gin-api/controllers/conn"
	"gin-api/controllers/opentracing"
	"gin-api/controllers/ping"
	"gin-api/middlewares/limiter"
	"gin-api/middlewares/log"
	"gin-api/middlewares/panic"
	"gin-api/middlewares/timeout"
	"gin-api/response"

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
		testGroup.POST("/rpc", opentracing.Rpc)
		testGroup.POST("/rpc1", opentracing.Rpc1)
		testGroup.POST("/panic", opentracing.Panic)
		testGroup.POST("/conn", conn.Do)
	}

	return server
}
