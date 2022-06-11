package rpc

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/why444216978/gin-api/library/logger"
)

// RPCConfig is used to parse configuration file
// logger should be controlled with Options
type RPCConfig struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

// RPCLogger is go-redis logger Hook
type RPCLogger struct {
	*logger.Logger
}

type RPCOption func(rl *RPCLogger)

// NewRPCLogger
func NewRPCLogger(cfg *RPCConfig, opts ...RPCOption) (rl *RPCLogger, err error) {
	rl = &RPCLogger{}

	for _, o := range opts {
		o(rl)
	}

	infoWriter, errWriter, err := logger.RotateWriter(cfg.InfoFile, cfg.ErrorFile)
	if err != nil {
		return
	}

	l, err := logger.NewLogger(&logger.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     cfg.Level,
	},
		logger.WithCallerSkip(5),
		logger.WithModule(logger.ModuleRPC),
		logger.WithInfoWriter(infoWriter),
		logger.WithErrorWriter(errWriter),
	)
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
	Request     interface{}
	Response    interface{}
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
	logFields := logger.ValueHTTPFields(ctx)

	logFields.Header = fields.Header
	logFields.Method = fields.Method
	logFields.ClientIP = logFields.ServerIP
	logFields.ClientPort = logFields.ServerPort
	logFields.ServerIP = fields.ServerIP
	logFields.ServerPort = fields.ServerPort
	logFields.API = fields.URI
	logFields.Request = fields.Request
	logFields.Response = fields.Response
	logFields.Cost = fields.Cost
	logFields.Code = fields.HTTPCode

	newCtx := context.WithValue(ctx, "rpc", "rpc")
	newCtx = logger.WithHTTPFields(newCtx, logFields)

	return newCtx, []zap.Field{
		zap.String(logger.SericeName, fields.ServiceName),
		zap.Duration(logger.Timeout, fields.Timeout),
	}
}
