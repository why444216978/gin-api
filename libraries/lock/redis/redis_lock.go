package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/why444216978/gin-api/libraries/lock"

	"github.com/go-redis/redis/v8"
)

const (
	lockSuccess = 1
	lockFail    = 0
	lockLua     = `if redis.call("GET", KEYS[1]) == ARGV[1] then redis.call("DEL", KEYS[1]) return 1 else return 0 end`
)

var _ lock.Locker = (*redisLock)(nil)

var (
	ErrClientNil = errors.New("client is nil")
	ErrLock      = errors.New("lock fail")
	ErrUnLock    = errors.New("unlock fail")
)

type redisLock struct {
	c *redis.Client
}

func New(c *redis.Client) (*redisLock, error) {
	if c == nil {
		return nil, ErrClientNil
	}
	return &redisLock{
		c: c,
	}, nil
}

// Lock lock
func (rl *redisLock) Lock(ctx context.Context, key, random string, duration time.Duration) (err error) {
	_, err = rl.c.SetNX(ctx, key, random, duration).Result()
	if err != nil {
		return
	}

	return
}

// UnLock unlock
func (rl *redisLock) Unlock(ctx context.Context, key, random string) (err error) {
	res, err := rl.c.Eval(ctx, lockLua, []string{key}, random).Result()
	if err != nil {
		return
	}

	if res == lockFail {
		err = errors.New(fmt.Sprintf("unlock result is %d", lockFail))
		return
	}

	return
}
