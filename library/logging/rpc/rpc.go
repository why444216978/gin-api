package logging

import (
	"context"
	"net/http"
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
	}, logging.WithCallerSkip(4), logging.WithModule(logging.ModuleRPC))
	if err != nil {
		return
	}
	rl.Logger = l

	return
}

type RPCLogFields struct {
	ServiceName string
	Header      http.Header
	Method      string
	URI         string
	Request     []byte
	Response    string
	ServerIP    string
	ServerPort  int
	HTTPCode    int
	Cost        int64
	Timeout     time.Duration
}

func (rl *RPCLogger) Info(ctx context.Context, msg string, fields RPCLogFields) {
	newCtx, logFields := rl.fields(ctx, fields)
	rl.Logger.Info(newCtx, msg, logFields...)
}

func (rl *RPCLogger) Error(ctx context.Context, msg string, fields RPCLogFields) {
	newCtx, logFields := rl.fields(ctx, fields)
	rl.Logger.Error(newCtx, msg, logFields...)
}

func (rl *RPCLogger) fields(ctx context.Context, fields RPCLogFields) (context.Context, []zap.Field) {
	//添加通用header
	fields.Header.Add(logging.LogHeader, logging.ValueLogID(ctx))

	response, _ := conversion.JsonToMap(fields.Response)
	request, _ := conversion.JsonToMap(string(fields.Request))

	logFields := logging.ValueHTTPFields(ctx)
	logFields.Method = fields.Method
	logFields.Header = fields.Header
	logFields.ClientIP = logFields.ServerIP
	logFields.ClientPort = logFields.ServerPort
	logFields.ServerIP = fields.ServerIP
	logFields.ServerPort = fields.ServerPort
	logFields.API = fields.URI
	logFields.Request = request
	logFields.Response = response
	logFields.Cost = fields.Cost
	logFields.Code = fields.HTTPCode

	newCtx := context.WithValue(ctx, "rpc", "rpc")
	newCtx = logging.WithHTTPFields(newCtx, logFields)

	return newCtx, []zap.Field{
		zap.String("service_name", fields.ServiceName),
		zap.Duration("timeout", fields.Timeout),
	}
}
