package rabbitmq

import (
	"fmt"

	"gin-api/libraries/logging"

	"github.com/streadway/amqp"
)

type TestConsumer struct{}

//消费mq消息
func (*TestConsumer) Do(d amqp.Delivery, header *logging.LogHeader) error {
	fmt.Println(string(d.Body))
	err := d.Ack(false)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
