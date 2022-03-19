package gin_api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jsonCodec "github.com/why444216978/codec/json"

	"github.com/why444216978/gin-api/app/resource"
	httpClient "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/logger"
)

const (
	serviceName = "gin-api"
)

func RPC(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:     fmt.Sprintf("/test/rpc1?logid=%s", logger.ValueLogID(ctx)),
		Method:  http.MethodPost,
		Header:  nil,
		Timeout: time.Second,
		Body:    map[string]interface{}{"rpc": "rpc"},
		Codec:   jsonCodec.JSONCodec{},
	}
	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, "gin-api-dev", req, resp); err != nil {
		return
	}

	return
}

func RPC1(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:     fmt.Sprintf("/test/conn?logid=%s", logger.ValueLogID(ctx)),
		Method:  http.MethodPost,
		Header:  nil,
		Timeout: time.Second,
		Body:    map[string]interface{}{"rpc1": "rpc1"},
		Codec:   jsonCodec.JSONCodec{},
	}

	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, "gin-api-dev", req, resp); err != nil {
		return
	}

	return
}

func Ping(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:     fmt.Sprintf("/ping?logid=%s", logger.ValueLogID(ctx)),
		Method:  http.MethodGet,
		Header:  nil,
		Timeout: time.Second,
		Body:    map[string]interface{}{"rpc1": "rpc1"},
		Codec:   jsonCodec.JSONCodec{},
	}

	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, "gin-api-dev", req, resp); err != nil {
		return
	}

	return
}
