package queue

import (
	"context"
)

type ProduceOption struct{}

type ProduceOptionFunc func(o *ProduceOption)

type ConsumeOption struct{}

type ConsumeOptionFunc func(o *ConsumeOption)

type Consumer func(context.Context, interface{}) (retry bool, err error)

type Queue interface {
	Produce(ctx context.Context, msg interface{}, opts ...ProduceOptionFunc) error
	Consume(consumer Consumer)
	Shutdown() error
}
