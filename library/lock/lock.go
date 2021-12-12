package lock

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrClientNil 客户端nil
	ErrClientNil = errors.New("client is nil")
	// ErrLock 加锁/获取锁失败
	ErrLock = errors.New("lock fail")
	// ErrUnLock 解锁失败
	ErrUnLock = errors.New("unlock fail")
)

type Locker interface {
	Lock(ctx context.Context, key string, random interface{}, duration time.Duration) (err error)
	Unlock(ctx context.Context, key string, random interface{}) (err error)
}
