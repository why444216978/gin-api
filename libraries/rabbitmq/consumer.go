package rabbitmq

import (
	"fmt"

	"gin-api/libraries/logging"

	"github.com/streadway/amqp"
)

type TestConsumer struct{}

//消费mq消息
func (self *TestConsumer) Do(d amqp.Delivery, header *logging.LogHeader) error {
	fmt.Println(string(d.Body))
	err := d.Ack(false)
	if err != nil {
		panic(err)
	}
	return err
}
