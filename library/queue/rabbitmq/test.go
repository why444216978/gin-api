package rabbitmq

import (
	"context"
	"log"
)

func Do(ctx context.Context, msg interface{}) (bool, error) {
	log.Printf("consume: %v", msg)
	return false, nil
}

func Test() {
	q := New("test", "amqp://why:why@localhost:5672/why", "why_exchange", ExchangeTypeDirect, "why_queue", "why_route")
	q.Produce(context.TODO(), []byte("test message"))

	q.Consume(Do)
}
