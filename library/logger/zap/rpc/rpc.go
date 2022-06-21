package rpc

import (
	"context"

	"github.com/why444216978/gin-api/library/logger"
	zapLogger "github.com/why444216978/gin-api/library/logger/zap"
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
	*zapLogger.ZapLogger
	config *RPCConfig
}

type RPCOption func(rl *RPCLogger)

// NewRPCLogger
func NewRPCLogger(config *RPCConfig, opts ...RPCOption) (rl *RPCLogger, err error) {
	rl = &RPCLogger{config: config}

	for _, o := range opts {
		o(rl)
	}

	infoWriter, errWriter, err := logger.RotateWriter(config.InfoFile, config.ErrorFile)
	if err != nil {
		return
	}

	l, err := zapLogger.NewLogger(
		zapLogger.WithCallerSkip(4),
		zapLogger.WithModule(logger.ModuleRPC),
		zapLogger.WithInfoWriter(infoWriter),
		zapLogger.WithErrorWriter(errWriter),
		zapLogger.WithLevel(config.Level),
	)
	if err != nil {
		return
	}
	rl.ZapLogger = l

	return
}

func (rl *RPCLogger) Info(ctx context.Context, msg string, fields logger.Fields) {
	newCtx, logFields := rl.fields(ctx, fields)
	rl.logger().Info(newCtx, msg, logFields...)
}

func (rl *RPCLogger) Error(ctx context.Context, msg string, fields logger.Fields) {
	newCtx, logFields := rl.fields(ctx, fields)
	rl.logger().Error(newCtx, msg, logFields...)
}

func (rl *RPCLogger) fields(ctx context.Context, fields logger.Fields) (context.Context, []logger.Field) {
	logFields := logger.ValueHTTPFields(ctx)

	logFields.Header = fields.Header
	logFields.Method = fields.Method
	logFields.ClientIP = logFields.ServerIP
	logFields.ClientPort = logFields.ServerPort
	logFields.ServerIP = fields.ServerIP
	logFields.ServerPort = fields.ServerPort
	logFields.API = fields.API
	logFields.Request = fields.Request
	logFields.Response = fields.Response
	logFields.Cost = fields.Cost
	logFields.Code = fields.Code

	newCtx := context.WithValue(ctx, "rpc", "rpc")
	newCtx = logger.WithHTTPFields(newCtx, logFields)

	return newCtx, []logger.Field{
		logger.Reflect(logger.ServiceName, fields.ServiceName),
		logger.Reflect(logger.Timeout, fields.Timeout),
	}
}

func (rl *RPCLogger) logger() *zapLogger.ZapLogger {
	return rl.ZapLogger
}
