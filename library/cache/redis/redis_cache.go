package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	util_ctx "github.com/why444216978/go-util/context"
	"github.com/why444216978/go-util/snowflake"

	"github.com/why444216978/gin-api/library/cache"
	"github.com/why444216978/gin-api/library/lock"
)

type RedisCache struct {
	c    *redis.Client
	lock lock.Locker
}

var _ cache.Cacher = (*RedisCache)(nil)

func New(c *redis.Client, locker lock.Locker) (*RedisCache, error) {
	if c == nil {
		return nil, errors.New("redis is nil")
	}

	if locker == nil {
		return nil, errors.New("locker is nil")
	}

	return &RedisCache{
		c:    c,
		lock: locker,
	}, nil
}

func (rc *RedisCache) GetData(ctx context.Context, key string, ttl time.Duration, virtualTTL time.Duration, f cache.LoadFunc, data interface{}) (err error) {
	cache, err := rc.getCache(ctx, key)
	if err != nil {
		return
	}

	// 无缓存
	if cache.ExpireAt == 0 || cache.Data == "" {
		err = rc.FlushCache(ctx, key, ttl, virtualTTL, f, data)
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
	go rc.FlushCache(ctxNew, key, ttl, virtualTTL, f, data)

	return
}

func (rc *RedisCache) FlushCache(ctx context.Context, key string, ttl time.Duration, virtualTTL time.Duration, f cache.LoadFunc, data interface{}) (err error) {
	lockKey := "LOCK::" + key
	random := snowflake.Generate().String()

	//获取锁，自旋三次
	//TODO 这里可优化为客户端传入控制
	try := 0
	for {
		try = try + 1
		if try > 3 {
			break
		}
		err = rc.lock.Lock(ctx, lockKey, random, time.Second*10)
		if err == lock.ErrLock {
			continue
		}
		break
	}
	if err != nil {
		return
	}
	defer rc.lock.Unlock(ctx, lockKey, random)

	// load data
	err = cache.HandleLoad(ctx, f, data)
	if err != nil {
		return
	}

	dataStr, err := json.Marshal(data)
	if err != nil {
		return
	}

	// save cache
	err = rc.setCache(ctx, key, string(dataStr), ttl, virtualTTL)

	return
}

func (rc *RedisCache) getCache(ctx context.Context, key string) (data *cache.CacheData, err error) {
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

func (rc *RedisCache) setCache(ctx context.Context, key, val string, ttl time.Duration, virtualTTL time.Duration) (err error) {
	_data := cache.CacheData{
		ExpireAt: time.Now().Add(virtualTTL).Unix(),
		Data:     val,
	}
	data, err := json.Marshal(_data)
	if err != nil {
		return
	}

	_, err = rc.c.Set(ctx, key, string(data), ttl).Result()
	if err != nil {
		return
	}

	return
}
