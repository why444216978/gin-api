package cache

import (
	"context"
	"time"
)

// cacheData cache data struct
type cacheData struct {
	ExpireAt int64  // ExpireAt 失效时间
	Data     string // Data 真实数据
}

// LoadFunc define load data func
type LoadFunc func(ctx context.Context, target interface{}) (err error)

type Cacher interface {
	// GetData load data from cache
	// if cache not exist load data by LoadFunc
	// expiration is redis server expiration
	// ttl is developer expiration
	GetData(ctx context.Context, key string, expiration time.Duration, ttl int64, f LoadFunc, data interface{}) (err error)

	// flushCache flush cache
	// if cache not exist, load data and save cache
	flushCache(ctx context.Context, key string, expiration time.Duration, ttl int64, f LoadFunc, data interface{}) (err error)

	// getCache get data from cache and conversion to cacheData
	getCache(ctx context.Context, key string) (data *cacheData, err error)

	// setCache set cache
	setCache(ctx context.Context, key, val string, expiration time.Duration, ttl int64) (err error)
}
