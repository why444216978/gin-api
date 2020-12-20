package routers

import (
	"gin-api/controllers/conn"
	"gin-api/controllers/opentracing"
	"gin-api/controllers/ping"
	"gin-api/middlewares/limiter"
	"gin-api/middlewares/log"
	"gin-api/middlewares/panic"
	"gin-api/middlewares/trace"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	server := gin.New()

	server.Use(gin.Recovery())

	server.Use(trace.OpenTracing())

	server.Use(limiter.Limiter(10))

	server.Use(log.LoggerMiddleware())

	server.Use(panic.ThrowPanic())

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
		testGroup.GET("/panic", opentracing.Panic)
		testGroup.GET("/conn", conn.Do)
	}

	return server
}
