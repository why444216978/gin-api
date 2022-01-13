package gin_api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/why444216978/gin-api/app/resource"
	httpClient "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/logger"
)

const (
	serviceName = "gin-api"
)

func RPC(ctx context.Context) (ret httpClient.Response, err error) {
	uri := fmt.Sprintf("/test/rpc1?logid=%s", logger.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.ClientHTTP.Send(ctx, "gin-api-dev", http.MethodPost, uri, nil, time.Second, map[string]interface{}{"rpc": "rpc"}, &resp)
}

func RPC1(ctx context.Context) (ret httpClient.Response, err error) {
	uri := fmt.Sprintf("/test/conn?logid=%s", logger.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.ClientHTTP.Send(ctx, "gin-api-dev", http.MethodPost, uri, nil, time.Second, map[string]interface{}{"rpc1": "rpc1"}, &resp)
}

func Ping(ctx context.Context) (ret httpClient.Response, err error) {
	uri := fmt.Sprintf("/ping?logid=%s", logger.ValueLogID(ctx))

	var resp map[string]interface{}
	return resource.ClientHTTP.Send(ctx, "gin-api-dev", http.MethodGet, uri, nil, time.Second, nil, &resp)
}
