package logging

import (
	"context"
	"time"

	"github.com/why444216978/gin-api/library/logging"

	"github.com/why444216978/go-util/conversion"
	"go.uber.org/zap"
)

type RPCConfig struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

// RPCLogger is go-redis logger Hook
type RPCLogger struct {
	*logging.Logger
}

type RPCOption func(rl *RPCLogger)

// NewRPCLogger
func NewRPCLogger(cfg *RPCConfig, opts ...RPCOption) (rl *RPCLogger, err error) {
	rl = &RPCLogger{}

	for _, o := range opts {
		o(rl)
	}

	l, err := logging.NewLogger(&logging.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     cfg.Level,
	}, logging.WithCallerSkip(2), logging.WithModule(logging.ModuleRPC))
	if err != nil {
		return
	}
	rl.Logger = l

	return
}

func (rl *RPCLogger) Fields(ctx context.Context, serviceName, method, uri string, header map[string]string, body []byte, timeout time.Duration,
	remoteHost string, remotePort int, resp string, err error) []zap.Field {

	response, _ := conversion.JsonToMap(resp)
	return []zap.Field{
		zap.String(logging.LogID, logging.ValueTraceID(ctx)),
		zap.String(logging.TraceID, logging.ValueLogID(ctx)),
		zap.String("service_name", serviceName),
		zap.String("method", method),
		zap.String("uri", uri),
		zap.String("remote_host", remoteHost),
		zap.Int("remote_port", remotePort),
		zap.Reflect("header", header),
		zap.Reflect("request", string(body)),
		zap.Reflect("response", response),
		zap.Error(err),
	}
}
