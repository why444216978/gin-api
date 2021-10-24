package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// ConsumeMsg
type ConsumeMsg interface {
	// Do
	Do(d amqp.Delivery) error
}

// Consumer
type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	amqpURI    string
	tag        string
	done       chan error // 业务逻辑结束标识
}

// NewConsumer
func NewConsumer(amqpURI, queueName, tag string, consumeMsg ConsumeMsg) (*Consumer, error) {
	c := &Consumer{
		amqpURI:   amqpURI,
		queueName: queueName,
		tag:       tag,
		done:      make(chan error),
	}
	var err error
	c.connection, err = amqp.Dial(c.amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial:%s", err)
	}
	defer c.connection.Close()

	go func() {
		fmt.Printf("closing: %s", <-c.connection.NotifyClose(make(chan *amqp.Error)))
	}()

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

	for d := range deliveries {
		go consume(d, consumeMsg)
	}
	<-c.done

	return c, nil
}

func consume(d amqp.Delivery, consumeMsg ConsumeMsg) (err error) {
	return consumeMsg.Do(d)
}

func NewProducer(ctx context.Context, msg, amqpURI, exchangeName, exchangeType, queueName, routeName, tag string) (err error) {
	c := &Consumer{
		amqpURI:   amqpURI,
		queueName: queueName,
		tag:       tag,
		done:      make(chan error),
	}
	c.connection, err = amqp.Dial(c.amqpURI)
	if err != nil {
		return fmt.Errorf("Dial:%s", err)
	}
	defer c.connection.Close()

	c.channel, err = c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	defer c.channel.Close()

	err = c.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		return
	}

	_, err = c.channel.QueueDeclare(
		queueName, // routing_key
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return
	}

	err = c.channel.QueueBind(queueName, routeName, exchangeName, false, nil)
	if err != nil {
		return
	}

	err = c.channel.Publish(
		exchangeName, // exchange
		routeName,    // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return
	}

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

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	<-c.done
}
