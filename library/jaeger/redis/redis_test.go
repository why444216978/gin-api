package redis

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"github.com/why444216978/gin-api/library/jaeger"
)

func TestNewJaegerHook(t *testing.T) {
	convey.Convey("TestNewJaegerHook", t, func() {
		convey.Convey("success", func() {
			_ = NewJaegerHook()
		})
	})
}

func Test_jaegerHook_BeforeProcess(t *testing.T) {
	convey.Convey("Test_jaegerHook_BeforeProcess", t, func() {
		convey.Convey("Tracer nil", func() {
			ctx := context.Background()
			jaeger.Tracer = nil
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcess(ctx, cmd)
			assert.Equal(t, err, nil)
		})
		convey.Convey("success", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcess(ctx, cmd)
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 0)
		})
	})
}

func Test_jaegerHook_AfterProcess(t *testing.T) {
	convey.Convey("Test_jaegerHook_AfterProcess", t, func() {
		convey.Convey("Tracer nil", func() {
			ctx := context.Background()
			jaeger.Tracer = nil
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			err := jh.AfterProcess(ctx, cmd)
			assert.Equal(t, err, nil)
		})
		convey.Convey("extract span from ctx nil", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			err := jh.AfterProcess(ctx, cmd)
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("success", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcess(ctx, cmd)
			err = jh.AfterProcess(ctx, cmd)
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 1)
		})
		convey.Convey("success and cmd err", func() {
			patche := gomonkey.ApplyFuncSeq(isRedisError, []gomonkey.OutputCell{
				{Values: gomonkey.Params{true}},
			})
			defer patche.Reset()

			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewBoolResult(false, redis.ErrClosed)
			ctx, err := jh.BeforeProcess(ctx, cmd)
			err = jh.AfterProcess(ctx, cmd)
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 1)
		})
	})
}

func Test_jaegerHook_BeforeProcessPipeline(t *testing.T) {
	convey.Convey("Test_jaegerHook_BeforeProcessPipeline", t, func() {
		convey.Convey("Tracer nil", func() {
			ctx := context.Background()
			jaeger.Tracer = nil
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
		})
		convey.Convey("success", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 0)
		})
	})
}

func Test_jaegerHook_AfterProcessPipeline(t *testing.T) {
	convey.Convey("Test_jaegerHook_AfterProcessPipeline", t, func() {
		convey.Convey("Tracer nil", func() {
			ctx := context.Background()
			jaeger.Tracer = nil
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			err := jh.AfterProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
		})
		convey.Convey("extract span from ctx nil", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			err := jh.AfterProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("success", func() {
			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewStringCmd(ctx, "get")
			ctx, err := jh.BeforeProcessPipeline(ctx, []redis.Cmder{cmd})
			err = jh.AfterProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 1)
		})
		convey.Convey("success and cmd err", func() {
			patche := gomonkey.ApplyFuncSeq(isRedisError, []gomonkey.OutputCell{
				{Values: gomonkey.Params{true}},
			})
			defer patche.Reset()

			ctx := context.Background()
			tracer := mocktracer.New()
			jaeger.Tracer = tracer
			jh := NewJaegerHook()
			cmd := redis.NewBoolResult(false, redis.ErrClosed)
			ctx, err := jh.BeforeProcessPipeline(ctx, []redis.Cmder{cmd})
			err = jh.AfterProcessPipeline(ctx, []redis.Cmder{cmd})
			assert.Equal(t, err, nil)
			assert.Len(t, tracer.FinishedSpans(), 1)
		})
	})
}

func Test_jaegerHook_getPipeLineLogKey(t *testing.T) {
	convey.Convey("Test_jaegerHook_getPipeLineLogKey", t, func() {
		convey.Convey("success", func() {
			assert.Equal(t, (&jaegerHook{}).getPipeLineLogKey("a", 1), "a-1")
		})
	})
}

func Test_isRedisError(t *testing.T) {
	convey.Convey("Test_isRedisError", t, func() {
		convey.Convey("redis.Nil", func() {
			assert.Equal(t, isRedisError(redis.Nil), false)
		})
		convey.Convey("not redis.Nil", func() {
			assert.Equal(t, isRedisError(nil), false)
		})
	})
}
