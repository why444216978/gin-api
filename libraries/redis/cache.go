package redis

import (
	"context"
	"encoding/json"
	"gin-api/libraries/logging"
	"net/http"
	"time"
)

// CacheData cache data struct
type CacheData struct {
	ExpireAt int64  // ExpireAt 失效时间
	Data     string // Data 真实数据
}

// LoadFunc define load data func
type LoadFunc func(ctx context.Context, target interface{}) (err error)

// GetData load data from cache
// if cache not exist load data by LoadFunc
func (db *RedisDB) GetData(ctx context.Context, header http.Header, key string, ttl, ex int64, f LoadFunc, data interface{}) (err error) {
	key = db.getCacheKey(key)

	cache, err := db.getCache(ctx, header, key)
	if err != nil {
		return
	}

	// 无缓存
	if cache.ExpireAt == 0 || cache.Data == "" {
		db.flushCache(ctx, header, key, ttl, ex, f, data)
		return
	}

	err = json.Unmarshal([]byte(cache.Data), data)
	if err != nil {
		return
	}

	if time.Now().Before(time.Unix(cache.ExpireAt, 0)) {
		return
	}

	go db.flushCache(ctx, header, key, ttl, ex, f, data)

	return
}

// flushCache flush cache
// if cache not exist, load data and save cache
func (db *RedisDB) flushCache(ctx context.Context, header http.Header, key string, ttl, ex int64, f LoadFunc, data interface{}) (err error) {
	lockKey := db.GetLockKey(key)
	uniqueStr := logging.NewObjectId().Hex()

	// lock
	err = db.Lock(ctx, header, lockKey, uniqueStr)
	if err != nil {
		return
	}
	defer db.UnLock(ctx, header, lockKey, uniqueStr)

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
	err = db.setCache(ctx, header, key, string(dataStr), ttl, ex)

	return
}

// getCache get data from cache
func (db *RedisDB) getCache(ctx context.Context, header http.Header, key string) (data *CacheData, err error) {
	data = &CacheData{}

	res, err := db.String(ctx, header, "GET", key)
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

// setCache set cache
func (db *RedisDB) setCache(ctx context.Context, header http.Header, key, val string, ttl, ex int64) (err error) {
	_data := CacheData{
		ExpireAt: time.Now().Unix() + ttl,
		Data:     val,
	}
	data, err := json.Marshal(_data)
	if err != nil {
		return
	}

	res, err := db.String(ctx, header, "SET", key, data, "ex", ex)
	if err != nil {
		err = ErrSetCache
		return
	}

	if res != ResultOK {
		err = ErrLock
		return
	}

	return
}

// getCacheKey get cache key
func (db *RedisDB) getCacheKey(key string) string {
	return "CACHE::" + key
}
