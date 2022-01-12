package http

import (
	"context"
	"net/http"

	jaeger "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logging"
)

type BeforeRequestPlugin interface {
	Handle(ctx context.Context, req *http.Request) error
}

type AfterRequestPlugin interface {
	Handle(ctx context.Context, req *http.Request, resp *http.Response) error
}

type JaegerBeforePlugin struct{}

var _ BeforeRequestPlugin = (*JaegerBeforePlugin)(nil)

func (*JaegerBeforePlugin) Handle(ctx context.Context, req *http.Request) error {
	jaeger.InjectHTTP(ctx, req, req.Header.Get(logging.LogHeader))

	return nil
}
