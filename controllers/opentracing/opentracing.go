package opentracing

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gin-api/libraries/config"
	"gin-api/response"

	"github.com/why444216978/go-util/conversion"
	util_http "github.com/why444216978/go-util/http"
)

func Rpc(c *gin.Context) {
	postData := make(map[string]interface{})
	postData["query"] = [1]string{"猕猴桃"}

	logCfg := config.GetConfigToJson("log", "log")
	logId := c.Writer.Header().Get(logCfg["query_field"].(string))
	sendUrl := "https://www.baidu.com?logid=" + logId

	body, _ := conversion.MapToJson(postData)

	ret, err := util_http.Send(c.Request.Context(), http.MethodGet, sendUrl, nil, strings.NewReader(body), time.Second)
	fmt.Println(ret)
	fmt.Println(err)

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
