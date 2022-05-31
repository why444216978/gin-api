package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/gin-api/server/http/response"
)

const (
	CodeSuccess     response.Code = 0
	CodeParams      response.Code = 1
	CodeUriNotFound response.Code = http.StatusNotFound
	CodeServer      response.Code = http.StatusInternalServerError
	CodeUnavailable response.Code = http.StatusServiceUnavailable
	CodeTimeout     response.Code = http.StatusGatewayTimeout
)

var codeToast = map[response.Code]string{
	CodeSuccess:     "success",
	CodeParams:      "参数错误",
	CodeUriNotFound: "资源不存在",
	CodeUnavailable: "服务器暂时不可用",
	CodeTimeout:     "请求超时",
	CodeServer:      "服务器错误",
}

func ResponseJSON(c *gin.Context, code response.Code, data interface{}, err *response.ResponseError) {
	if err == nil {
		err = response.WrapToast(nil, codeToast[code])
	}

	response.ResponseJSON(c, code, data, err)
}
