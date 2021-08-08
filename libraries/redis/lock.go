package redis

import (
	"context"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

const (
	lockSuccess = 1
	lockFail    = 0
	lockLua     = `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					redis.call("DEL", KEYS[1])
					return 1
				else
					return 0
				end
				`
)

// Lock lock
func (db *RedisDB) Lock(ctx context.Context, header http.Header, key, uniqueStr string) (err error) {
	res, err := db.String(ctx, header, "SET", key, uniqueStr, "ex", 10, "nx")
	if err != nil {
		return
	}

	if res != ResultOK {
		err = ErrLock
		return
	}

	return
}

// UnLock unlock
func (db *RedisDB) UnLock(ctx context.Context, header http.Header, key, uniqueStr string) (err error) {
	res, err := redis.Int(db.DoLua(ctx, header, lockLua, key, uniqueStr))
	if err != nil || res == lockFail {
		err = ErrUnLock
		return
	}

	return
}

// GetLockKey get lock key
func (db *RedisDB) GetLockKey(key string) string {
	return "LOCK_KEY::" + key
}
