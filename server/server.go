package server

import (
	"context"
	"errors"
)

type Server interface {
	Start() error
	Close() error
}

// CloseFunc 资源回收方法列表
var CloseFunc = make([]func(ctx context.Context) error, 0)

// RegisterCloseFunc 注册资源回收方法
func RegisterCloseFunc(cf interface{}) error {
	f, ok := cf.(func(ctx context.Context) error)
	if !ok {
		return errors.New("func type error")
	}

	CloseFunc = append(CloseFunc, f)
	return nil
}
