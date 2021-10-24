package gin_api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	lib_http "github.com/why444216978/gin-api/libraries/http"
	"github.com/why444216978/gin-api/libraries/logging"
	"github.com/why444216978/gin-api/resource"
)

const (
	serviceName = "gin-api"
)

func RPC(ctx context.Context) (ret *lib_http.Response, err error) {
	uri := fmt.Sprintf("/test/rpc1?logid=%s", logging.ValueLogID(ctx))
	header := map[string]string{logging.LogHeader: logging.ValueLogID(ctx)}

	return resource.HTTPRPC.Send(ctx, "gin-api", http.MethodPost, uri, header, bytes.NewBufferString(`{"a":"a"}`), time.Second)
}

func RPC1(ctx context.Context) (ret *lib_http.Response, err error) {
	uri := fmt.Sprintf("/test/conn?logid=%s", logging.ValueLogID(ctx))
	header := map[string]string{logging.LogHeader: logging.ValueLogID(ctx)}

	return resource.HTTPRPC.Send(ctx, "gin-api", http.MethodPost, uri, header, nil, time.Second)
}
