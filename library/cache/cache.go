package cache

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

// CacheData is cache data struct
type CacheData struct {
	ExpireAt int64  // ExpireAt is virtual expire time
	Data     string // Data is cache data
}

// LoadFunc is define load data func
type LoadFunc func(ctx context.Context, target interface{}) (err error)

// Cacher is used to load cache
type Cacher interface {
	// GetData load data from cache
	// if cache not exist load data by LoadFunc
	// ttl is redis server ttl
	// virtualTTL is developer ttl
	GetData(ctx context.Context, key string, ttl time.Duration, virtualTTL time.Duration, f LoadFunc, data interface{}) (err error)

	// FlushCache flush cache
	// if cache not exist, load data and save cache
	FlushCache(ctx context.Context, key string, ttl time.Duration, virtualTTL time.Duration, f LoadFunc, data interface{}) (err error)
}

// panicError https://cs.opensource.google/go/x/sync/+/036812b2:singleflight/singleflight.go
type panicError struct {
	value interface{}
	stack []byte
}

// Error implements error interface.
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

// newPanicError is format panic error
func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// HandleLoad is used load cache
func HandleLoad(ctx context.Context, f LoadFunc, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newPanicError(r)
		}
	}()
	err = f(ctx, data)
	return
}
