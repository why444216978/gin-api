package routers

import (
	"gin-frame/controllers/first_origin_price"
	"gin-frame/controllers/opentracing"
	"gin-frame/controllers/ping"
	"gin-frame/libraries/config"
	"gin-frame/middlewares/limiter"
	"gin-frame/middlewares/log"
	"gin-frame/middlewares/panic"
	"gin-frame/middlewares/trace"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(port int, productName, moduleName, env string) *gin.Engine {
	logFields := config.GetLogFields()

	server := gin.New()

	server.Use(gin.Recovery())

	server.Use(trace.OpenTracing(productName))

	server.Use(limiter.Limiter(10))

	server.Use(log.LoggerMiddleware(port, logFields, productName, moduleName, env))

	server.Use(panic.ThrowPanic(port, logFields, productName, moduleName, env))
	//server.Use(dump.BodyDump())

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
	}

	originGroup := server.Group("/origin")
	{
		originGroup.GET("/first_origin_price", first_origin_price.Do)
	}

	return server
}
