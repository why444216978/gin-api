package redis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/why444216978/gin-api/library/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisConfig struct {
	InfoFile    string
	ErrorFile   string
	Level       string
	ServiceName string
	Host        string
	Port        int
}

// RedisLogger is go-redis logger Hook
type RedisLogger struct {
	*logger.Logger
	Config RedisConfig
}

type RedisOption func(rl *RedisLogger)

// NewRedisLogger
func NewRedisLogger(cfg *RedisConfig, opts ...RedisOption) (rl *RedisLogger, err error) {
	rl = &RedisLogger{
		Config: *cfg,
	}

	for _, o := range opts {
		o(rl)
	}

	l, err := logger.NewLogger(&logger.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     cfg.Level,
	}, logger.WithModule(logger.ModuleRedis), logger.WithServiceName(logger.ModuleRedis), logger.WithCallerSkip(5))
	if err != nil {
		return
	}
	rl.Logger = l

	return
}

//BeforeProcess redis before execute action do something
func (rl *RedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = rl.setCmdStart(ctx)
	return ctx, nil
}

//AfterProcess redis after execute action do something
func (rl *RedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if rl.Logger == nil {
		return nil
	}

	cost := rl.getCmdCost(ctx)

	if err := cmd.Err(); err != nil && err != redis.Nil {
		rl.Error(ctx, cmd, cost)
		return nil
	}

	rl.Info(ctx, cmd, cost)

	return nil
}

// BeforeProcessPipeline before command process handle
func (rl *RedisLogger) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

// AfterProcessPipeline after command process handle
func (rl *RedisLogger) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if rl.Logger == nil {
		return nil
	}

	for idx, cmd := range cmds {
		cost := rl.getPiplineCost(ctx, idx)

		if err := cmd.Err(); err != nil && err != redis.Nil {
			rl.Error(ctx, cmd, cost)
			continue
		}
		rl.Info(ctx, cmd, cost)
	}

	return nil
}

func (rl *RedisLogger) Info(ctx context.Context, cmd redis.Cmder, cost int64) {
	newCtx, logFields := rl.fields(ctx, cmd, cost)
	rl.Logger.Info(newCtx, "info", logFields...)
}

func (rl *RedisLogger) Error(ctx context.Context, cmd redis.Cmder, cost int64) {
	newCtx, logFields := rl.fields(ctx, cmd, cost)
	rl.Logger.Error(newCtx, cmd.Err().Error(), logFields...)
}

func (rl *RedisLogger) fields(ctx context.Context, cmd redis.Cmder, cost int64) (context.Context, []zap.Field) {
	logFields := logger.ValueHTTPFields(ctx)
	logFields.Header = http.Header{}
	logFields.Method = cmd.Name()
	logFields.Request = cmd.Args()
	logFields.Response = cmd.String()
	logFields.Code = 0
	logFields.ClientIP = logFields.ServerIP
	logFields.ClientPort = logFields.ServerPort
	logFields.ServerIP = rl.Config.Host
	logFields.ServerPort = rl.Config.Port
	logFields.API = cmd.Name()
	logFields.Cost = cost

	newCtx := context.WithValue(ctx, "rpc", "rpc")
	newCtx = logger.WithHTTPFields(newCtx, logFields)
	return newCtx, []zap.Field{
		zap.String("service_name", rl.Config.ServiceName),
	}
}

const (
	cmdStart        = "cmdStart"
	piplineCmdStart = "piplineCmdStart_"
)

func (rl *RedisLogger) setCmdStart(ctx context.Context) context.Context {
	return context.WithValue(ctx, cmdStart, time.Now())
}

func (rl *RedisLogger) getCmdCost(ctx context.Context) int64 {
	return time.Since(ctx.Value(cmdStart).(time.Time)).Milliseconds()
}

func (rl *RedisLogger) setPiplineStart(ctx context.Context, idx int) context.Context {
	return context.WithValue(ctx, fmt.Sprintf("%s%d", piplineCmdStart, idx), time.Now())
}

func (rl *RedisLogger) getPiplineCost(ctx context.Context, idx int) int64 {
	return time.Since(ctx.Value(fmt.Sprintf("%s%d", piplineCmdStart, idx)).(time.Time)).Milliseconds()
}
