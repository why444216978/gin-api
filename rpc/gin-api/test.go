package gin_api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/why444216978/gin-api/library/logging"
	lib_http "github.com/why444216978/gin-api/library/rpc/http"
	"github.com/why444216978/gin-api/resource"
)

const (
	serviceName = "gin-api"
)

func RPC(ctx context.Context) (ret lib_http.Response, err error) {
	uri := fmt.Sprintf("/test/rpc1?logid=%s", logging.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.HTTPRPC.Send(ctx, "gin-api-dev", http.MethodPost, uri, nil, time.Second, map[string]interface{}{"rpc": "rpc"}, &resp)
}

func RPC1(ctx context.Context) (ret lib_http.Response, err error) {
	uri := fmt.Sprintf("/test/conn?logid=%s", logging.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.HTTPRPC.Send(ctx, "gin-api-dev", http.MethodPost, uri, nil, time.Second, map[string]interface{}{"rpc1": "rpc1"}, &resp)
}

func Ping(ctx context.Context) (ret lib_http.Response, err error) {
	uri := fmt.Sprintf("/ping?logid=%s", logging.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.HTTPRPC.Send(ctx, "gin-api-dev", http.MethodGet, uri, nil, time.Second, nil, &resp)
}
