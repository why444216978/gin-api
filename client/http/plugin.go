package http

import (
	"context"
	"net/http"

	jaeger "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logger"
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
	logID := logger.ValueLogID(ctx)
	req.Header.Add(logger.LogHeader, logID)
	return jaeger.InjectHTTP(ctx, req, logID)
}
