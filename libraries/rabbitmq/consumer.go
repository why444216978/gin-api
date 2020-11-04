package rabbitmq

import (
	"fmt"

	"gin-api/libraries/log"
	"gin-api/libraries/util"

	"github.com/streadway/amqp"
)

type TestConsumer struct{}

//消费mq消息
func (self *TestConsumer) Do(d amqp.Delivery, header *log.LogFormat) error {
	fmt.Println(string(d.Body))
	err := d.Ack(false)
	util.Must(err)
	return err
}
