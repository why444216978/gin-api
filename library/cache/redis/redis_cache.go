package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/why444216978/gin-api/library/cache"
	"github.com/why444216978/gin-api/library/lock"
	"github.com/why444216978/gin-api/library/logging"

	"github.com/go-redis/redis/v8"
	util_ctx "github.com/why444216978/go-util/context"
)

type redisCache struct {
	c    *redis.Client
	lock lock.Locker
}

var _ cache.Cacher = (*redisCache)(nil)

func New(c *redis.Client, locker lock.Locker) (*redisCache, error) {
	if c == nil {
		return nil, errors.New("redis is nil")
	}

	if locker == nil {
		return nil, errors.New("locker is nil")
	}

	return &redisCache{
		c:    c,
		lock: locker,
	}, nil
}

func (rc *redisCache) GetData(ctx context.Context, key string, expiration time.Duration, ttl time.Duration, f cache.LoadFunc, data interface{}) (err error) {
	cache, err := rc.getCache(ctx, key)
	if err != nil {
		return
	}

	// 无缓存
	if cache.ExpireAt == 0 || cache.Data == "" {
		rc.FlushCache(ctx, key, expiration, ttl, f, data)
		return
	}

	err = json.Unmarshal([]byte(cache.Data), data)
	if err != nil {
		return
	}

	if time.Now().Before(time.Unix(cache.ExpireAt, 0)) {
		return
	}

	ctxNew := util_ctx.RemoveCancel(ctx)
	go rc.FlushCache(ctxNew, key, expiration, ttl, f, data)

	return
}

func (rc *redisCache) FlushCache(ctx context.Context, key string, expiration time.Duration, ttl time.Duration, f cache.LoadFunc, data interface{}) (err error) {
	lockKey := "LOCK::" + key
	random := logging.NewObjectId().Hex()

	// lock
	err = rc.lock.Lock(ctx, lockKey, random, time.Second*10)
	if err != nil {
		return
	}
	defer rc.lock.Unlock(ctx, lockKey, random)

	// load data
	err = f(ctx, data)
	if err != nil {
		return
	}

	dataStr, err := json.Marshal(data)
	if err != nil {
		return
	}

	// save cache
	err = rc.setCache(ctx, key, string(dataStr), expiration, ttl)

	return
}

func (rc *redisCache) getCache(ctx context.Context, key string) (data *cache.CacheData, err error) {
	data = &cache.CacheData{}

	res, err := rc.c.Get(ctx, key).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		return
	}

	if res == "" {
		return
	}

	err = json.Unmarshal([]byte(res), data)
	if err != nil {
		return
	}

	return
}

func (rc *redisCache) setCache(ctx context.Context, key, val string, expiration time.Duration, ttl time.Duration) (err error) {
	_data := cache.CacheData{
		ExpireAt: time.Now().Add(ttl).Unix(),
		Data:     val,
	}
	data, err := json.Marshal(_data)
	if err != nil {
		return
	}

	_, err = rc.c.Set(ctx, key, string(data), expiration).Result()
	if err != nil {
		return
	}

	return
}
