package router

import (
	"github.com/gin-gonic/gin"

	conn "github.com/why444216978/gin-api/app/module/goods/api"
	ping "github.com/why444216978/gin-api/app/module/ping/api"
	test "github.com/why444216978/gin-api/app/module/test/api"
)

func RegisterRouter(server *gin.Engine) {
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
}
