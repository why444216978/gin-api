package opentracing

import (
	"github.com/gin-gonic/gin"

	"gin-api/libraries/config"
	rpc_http "gin-api/libraries/http"
	"gin-api/response"
)

func Rpc(c *gin.Context) {
	postData := make(map[string]interface{})
	postData["query"] = [1]string{"猕猴桃"}

	logCfg := config.GetConfigToJson("log", "log")
	logId := c.Writer.Header().Get(logCfg["query_field"].(string))
	sendUrl := "https://www.baidu.com"

	ret := rpc_http.HttpSend(c, "GET", sendUrl, logId, postData)

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
