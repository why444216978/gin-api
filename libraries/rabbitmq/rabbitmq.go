package rabbitmq

import (
	"context"
	"fmt"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
	"github.com/streadway/amqp"
	srcLog "log"
	"net/http"
)

// 定义接收者接口
type ConsumeMsg interface {
	// 消费逻辑
	Do(d amqp.Delivery, header *log.LogFormat) error
}

// 定义RabbitMQ对象
type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string //队列名
	amqpURI    string
	tag        string
	done       chan error // 业务逻辑结束标识
}

func NewConsumer(amqpURI, queueName, tag string, consumeMsg ConsumeMsg) (*Consumer, error) {
	c := &Consumer{
		amqpURI:   amqpURI,
		queueName: queueName,
		tag:       tag,
		done:      make(chan error),
	}
	var err error
	srcLog.Printf("dialing %q", c.amqpURI)
	c.connection, err = amqp.Dial(c.amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial:%s", err)
	}
	defer c.connection.Close()

	go func() {
		fmt.Printf("closing: %s", <-c.connection.NotifyClose(make(chan *amqp.Error)))
	}()

	srcLog.Printf("got Connection, getting Channel")
	c.channel, err = c.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}
	defer c.channel.Close()

	deliveries, err := c.channel.Consume(
		queueName, // name
		c.tag,     // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	srcLog.Printf("consumer:(%q) started", queueName)

	for d := range deliveries {
		go consume(d, consumeMsg)
	}
	<-c.done

	return c, nil
}

func consume(d amqp.Delivery, consumeMsg ConsumeMsg) {
	logId := log.NewObjectId().Hex()

	logFormat := &log.LogFormat{
		LogId:    logId,
		CallerIp: "",
		HostIp:   "",
		Port:     0,
		Product:  "",
		Module:   "",
		//ServiceId:  serviceId,
		//InstanceId: host,
		UriPath: "",
		//XHop:    nil,
		//Tag:        "",
		Env: "",
		//SVersion:   "",
		//Stag: app.Stag{},
	}

	err := consumeMsg.Do(d, logFormat)
	util.Must(err)

	if err != nil {
		log.Errorf(logFormat, "failed to consumer msg:%s, err%s", d.Body, err.Error())
		util.Must(err)
		return
	}
}

func NewProducer(msg, amqpURI, exchangeName, exchangeType, queueName, routeName, tag string, ctx context.Context, header http.Header) error {
	var err error
	c := &Consumer{
		amqpURI:   amqpURI,
		queueName: queueName,
		tag:       tag,
		done:      make(chan error),
	}
	srcLog.Printf("dialing %q", c.amqpURI)
	c.connection, err = amqp.Dial(c.amqpURI)
	if err != nil {
		return fmt.Errorf("Dial:%s", err)
	}
	defer c.connection.Close()

	srcLog.Printf("got Connection, getting Channel")
	c.channel, err = c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	defer c.channel.Close()

	err = c.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	util.Must(err)

	_, err = c.channel.QueueDeclare(
		queueName, // routing_key
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	util.Must(err)

	err = c.channel.QueueBind(queueName, routeName, exchangeName, false, nil)
	util.Must(err)
	err = c.channel.Publish(
		exchangeName, // exchange
		routeName,    // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	util.Must(err)

	return nil
}

func NewProducerCmd(msg, amqpURI, exchangeName, exchangeType, queueName, routeName, tag string) error {
	var err error

	c := &Consumer{
		amqpURI:   amqpURI,
		queueName: queueName,
		tag:       tag,
		done:      make(chan error),
	}
	srcLog.Printf("dialing %q", c.amqpURI)

	c.connection, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer c.connection.Close()

	srcLog.Printf("got Connection, getting Channel")
	c.channel, err = c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	srcLog.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchangeName)
	if err = c.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	_, err = c.channel.QueueDeclare(
		queueName, // routing_key
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	util.Must(err)

	err = c.channel.QueueBind(queueName, routeName, exchangeName, false, nil)
	util.Must(err)

	err = c.channel.Publish(
		exchangeName, // exchange
		routeName,    // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	util.Must(err)

	return nil
}

func (c *Consumer) Shutdown() {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		fmt.Print(fmt.Sprintf("Consumer cancel failed: %s", err.Error()))
		return
	}

	if err := c.connection.Close(); err != nil {
		fmt.Print(fmt.Sprintf("AMQP connection close error: %s", err.Error()))
		return
	}

	defer srcLog.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	<-c.done
}
