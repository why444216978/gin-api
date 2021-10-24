package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type TestConsumer struct{}

func (*TestConsumer) Do(d amqp.Delivery) error {
	fmt.Println(string(d.Body))
	err := d.Ack(false)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func consumer() {
	connUri := "amqp://why:why@localhost:5672/why"
	queueName := "why_queue"

	testConsumer := &TestConsumer{}
	c, err := NewConsumer(connUri, queueName, "", testConsumer)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer c.Shutdown()
}

func producer() {
	connUri := "amqp://why:why@localhost:5672/why"
	queueName := "why_queue"
	exchangeName := "why_exchange"
	routeName := "why_route"

	err := NewProducer(context.TODO(), "hello world!", connUri, exchangeName, "direct", queueName, routeName, "")
	if err != nil {
		log.Println(err.Error())
	}
}

func Test() {
	producer()
	consumer()
}
