package lock

import (
	"context"
	"time"

	"github.com/why444216978/gin-api/library/lock"

	"github.com/go-redis/redis/v8"
)

const (
	lockSuccess = 1
	lockFail    = 0
	lockLua     = `if redis.call("GET", KEYS[1]) == ARGV[1] then redis.call("DEL", KEYS[1]) return 1 else return 0 end`
)

var _ lock.Locker = (*RedisLock)(nil)

type RedisLock struct {
	c *redis.Client
}

func New(c *redis.Client) (*RedisLock, error) {
	if c == nil {
		return nil, lock.ErrClientNil
	}
	return &RedisLock{
		c: c,
	}, nil
}

// Lock lock
func (rl *RedisLock) Lock(ctx context.Context, key string, random interface{}, duration time.Duration) (err error) {
	isSuccess, err := rl.c.SetNX(ctx, key, random, duration).Result()
	if err != nil {
		return
	}
	if !isSuccess {
		return lock.ErrLock
	}

	return
}

// UnLock unlock
func (rl *RedisLock) Unlock(ctx context.Context, key string, random interface{}) (err error) {
	res, err := rl.c.Eval(ctx, lockLua, []string{key}, random).Result()
	if err != nil {
		return
	}

	if res == lockFail {
		err = lock.ErrUnLock
		return
	}

	return
}
