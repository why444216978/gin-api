package routers

import (
	"gin-api/controllers/conn"
	"gin-api/controllers/opentracing"
	"gin-api/controllers/ping"
	"gin-api/middlewares/limiter"
	"gin-api/middlewares/log"
	"gin-api/middlewares/panic"
	"gin-api/middlewares/timeout"
	"gin-api/middlewares/trace"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	server := gin.New()

	server.Use(gin.Recovery())

	server.Use(log.WithContext())

	server.Use(panic.ThrowPanic())

	server.Use(limiter.Limiter(10))

	server.Use(timeout.TimeoutMiddleware(time.Second * 3))

	server.Use(log.LoggerMiddleware())

	server.Use(trace.OpenTracing())

	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"errno":    http.StatusNotFound,
			"errmsg":   "uri错误",
			"data":     nil,
			"user_msg": "请求资源不存在",
		})
	})

	pingGroup := server.Group("/ping")
	{
		pingGroup.GET("", ping.Ping)
	}

	testGroup := server.Group("/test")
	{
		testGroup.GET("/rpc", opentracing.Rpc)
		testGroup.POST("/panic", opentracing.Panic)
		testGroup.POST("/conn", conn.Do)
	}

	return server
}
