package ping

import (
	"gin-frame/libraries/config"
	"net/http"

	rpc_http "gin-frame/libraries/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	postData := make(map[string]interface{})
	postData["query"] = [1]string{"猕猴桃"}

	logId := c.Writer.Header().Get(config.GetHeaderLogIdField())
	sendUrl := "http://127.0.0.1:777/test/rpc"

	ret := rpc_http.HttpSend(c, "GET", sendUrl, logId, postData)
	ret = rpc_http.HttpSend(c, "GET", sendUrl, logId, postData)

	c.JSON(http.StatusOK, gin.H{
		"errno":  0,
		"errmsg": "success",
		"data":   ret,
	})
	c.Done()
}
