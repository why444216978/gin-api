package jaeger

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
)

const (
	ctxKey                         string = "redis_span"
	redisCmdName                   string = "command"
	redisCmdArgs                   string = "args"
	redisCmdResult                 string = "result"
	ErrNumberOfConnectionsExceeded string = "ERR max number of clients reached"
)

// jaegerHook is go-redis jaeger hook
type jaegerHook struct{}

// NewJaegerHook return jaegerHook
func NewJaegerHook() *jaegerHook {
	return &jaegerHook{}
}

//BeforeProcess redis before execute action do something
func (jh *jaegerHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if Tracer == nil {
		return ctx, nil
	}
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, Tracer, cmd.Name())
	if span == nil {
		return ctx, nil
	}
	span.SetTag(redisCmdName, cmd.Name())
	ctx = context.WithValue(ctx, ctxKey, span)
	return ctx, nil
}

//AfterProcess redis after execute action do something
func (jh *jaegerHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if Tracer == nil {
		return nil
	}
	_span := ctx.Value(ctxKey)
	span, ok := _span.(opentracing.Span)
	if !ok {
		return nil
	}
	defer span.Finish()

	if err := cmd.Err(); isRedisError(err) {
		span.LogFields(tracerLog.Error(err))
		span.SetTag(string(ext.Error), true)
	}
	span.LogFields(tracerLog.String(redisCmdName, cmd.Name()))
	span.LogFields(tracerLog.Object(redisCmdArgs, cmd.Args()))
	span.LogFields(tracerLog.String(redisCmdResult, cmd.String()))

	return nil
}

// BeforeProcessPipeline before command process handle
func (jh *jaegerHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if Tracer == nil {
		return ctx, nil
	}
	for _, cmd := range cmds {
		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, Tracer, cmd.Name())
		span.SetTag(redisCmdName, cmd.Name())
		ctx = context.WithValue(ctx, ctxKey, span)
	}
	return ctx, nil
}

// AfterProcessPipeline after command process handle
func (jh *jaegerHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if Tracer == nil {
		return nil
	}
	for _, cmd := range cmds {
		_span := ctx.Value(ctxKey)
		span, ok := _span.(opentracing.Span)
		if !ok {
			return nil
		}
		if err := cmd.Err(); isRedisError(err) {
			span.LogFields(tracerLog.Error(err))
			span.SetTag(string(ext.Error), true)
		}
		span.LogFields(tracerLog.String(redisCmdName, cmd.Name()))
		span.LogFields(tracerLog.Object(redisCmdArgs, cmd.Args()))
		span.LogFields(tracerLog.String(redisCmdResult, cmd.String()))
		span.Finish()
	}
	return nil
}

// redisError interface
type redisError interface {
	error

	// RedisError is a no-op function but
	// serves to distinguish types that are Redis
	// errors from ordinary errors: a type is a
	// Redis error if it has a RedisError method.
	RedisError()
}

func isRedisError(err error) bool {
	if err == redis.Nil {
		return false
	}
	_, ok := err.(redisError)
	return ok
}
