package lock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key, random string, duration time.Duration) (err error)
	Unlock(ctx context.Context, key, random string) (err error)
}
