package redis

import "errors"

var (
	ErrDataAssert = errors.New("data assert fail")
	ErrLock       = errors.New("lock fail")
	ErrUnLock     = errors.New("unlock fail")
	ErrSetCache   = errors.New("set cache fail")
)
