package logging

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/why444216978/gin-api/library/logging"
	"go.uber.org/zap"
)

type RedisConfig struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

// RedisLogger is go-redis logger Hook
type RedisLogger struct {
	*logging.Logger
}

type RedisOption func(rl *RedisLogger)

// NewRedisLogger
func NewRedisLogger(cfg *RedisConfig, opts ...RedisOption) (rl *RedisLogger, err error) {
	rl = &RedisLogger{}

	for _, o := range opts {
		o(rl)
	}

	l, err := logging.NewLogger(&logging.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     cfg.Level,
	}, logging.WithModule(logging.ModuleRedis))
	if err != nil {
		return
	}
	rl.Logger = l

	return
}

//BeforeProcess redis before execute action do something
func (rl *RedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

//AfterProcess redis after execute action do something
func (rl *RedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if rl.Logger == nil {
		return nil
	}

	var err error
	if e := cmd.Err(); e != redis.Nil {
		err = e
	}

	if err != nil {
		rl.Error("redis", rl.fields(ctx, cmd, err)...)
		return nil
	}
	rl.Info("redis", rl.fields(ctx, cmd, err)...)

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

	var err error

	for _, cmd := range cmds {
		if e := cmd.Err(); e != redis.Nil {
			err = e
		}
		if err != nil {
			rl.Logger.Error("redis", rl.fields(ctx, cmd, err)...)
			continue
		}
		rl.Logger.Info("redis", rl.fields(ctx, cmd, err)...)
	}

	return nil
}

func (rl *RedisLogger) fields(ctx context.Context, cmd redis.Cmder, err error) []zap.Field {
	return []zap.Field{
		zap.String(logging.LogID, logging.ValueTraceID(ctx)),
		zap.String(logging.TraceID, logging.ValueLogID(ctx)),
		zap.String("cmd", cmd.Name()),
		zap.Reflect("args", cmd.Args()),
		zap.String("result", cmd.String()),
		zap.Error(err),
	}
}
