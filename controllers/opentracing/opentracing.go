package opentracing

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-api/libraries/config"
	rpc_http "gin-api/libraries/http"
)

func Rpc(c *gin.Context) {
	postData := make(map[string]interface{})
	postData["query"] = [1]string{"猕猴桃"}

	logCfg := config.GetConfigToJson("log", "log")
	logId := c.Writer.Header().Get(logCfg["query_field"].(string))
	sendUrl := "https://www.baidu.com"

	ret := rpc_http.HttpSend(c, "GET", sendUrl, logId, postData)
	ret = rpc_http.HttpSend(c, "GET", sendUrl, logId, postData)

	c.JSON(http.StatusOK, gin.H{
		"errno":  0,
		"errmsg": "success",
		"data":   ret,
	})
	c.Done()
}

func Panic(c *gin.Context) {
	panic("test err")

	c.JSON(http.StatusOK, gin.H{
		"errno":  0,
		"errmsg": "success",
		"data":  make(map[string]interface{}),
	})
	c.Done()
}
