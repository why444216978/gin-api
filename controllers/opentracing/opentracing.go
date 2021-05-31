package opentracing

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gin-api/app_const"
	"gin-api/libraries/logging"
	"gin-api/response"

	"github.com/why444216978/go-util/conversion"
	util_http "github.com/why444216978/go-util/http"
)

func Rpc(c *gin.Context) {
	postData := make(map[string]interface{})
	postData["query"] = [1]string{"猕猴桃"}

	sendUrl := fmt.Sprintf("http://localhost:%d/test/conn?logid=%s", app_const.SERVICE_PORT, logging.GetLogID(c))

	body, _ := conversion.MapToJson(postData)

	ret, err := util_http.Send(c, http.MethodPost, sendUrl, nil, strings.NewReader(body), time.Second)
	fmt.Println(ret)
	fmt.Println(err)

	response.Response(c, response.CODE_SUCCESS, ret, "")
}

func Panic(c *gin.Context) {
	panic("test err")
}
