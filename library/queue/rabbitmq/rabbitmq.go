package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/why444216978/gin-api/library/queue"
)

const (
	ExchangeTypeDirect  = "direct"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeHeaders = "headers"
)

type Option struct {
	declareExchange bool
	declareQueue    bool
	bindQueue       bool
}

type OptionFunc func(*Option)

func defaultOption() *Option {
	return &Option{
		declareExchange: false,
		declareQueue:    false,
		bindQueue:       false,
	}
}

func WithDeclareExchange(turn bool) OptionFunc {
	return func(o *Option) { o.declareExchange = turn }
}

func WithDeclareQueue(turn bool) OptionFunc {
	return func(o *Option) { o.declareQueue = turn }
}

func WithBindQueue(turn bool) OptionFunc {
	return func(o *Option) { o.bindQueue = turn }
}

type RabbitMQ struct {
	opts         *Option
	connection   *amqp.Connection
	channel      *amqp.Channel
	name         string
	url          string
	exchangeName string
	exchangeType string
	queueName    string
	routeName    string
}

func New(name, url, exchangeName, exchangeType, queueName, routeName string, opts ...OptionFunc) *RabbitMQ {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}

	return &RabbitMQ{
		name:         name,
		url:          url,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
		routeName:    routeName,
	}
}

func (q *RabbitMQ) Produce(ctx context.Context, msg interface{}, opts ...queue.ProduceOptionFunc) (err error) {
	m, ok := msg.([]byte)
	if !ok {
		return errors.New("RabbitMQ msg not []byte")
	}

	if err = q.connect(); err != nil {
		return
	}

	err = q.channel.Publish(
		q.exchangeName, // exchange
		q.routeName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         m,
		})
	if err != nil {
		return errors.Wrap(err, "channel.Publish fail")
	}

	return
}

func (q *RabbitMQ) Consume(consumer queue.Consumer) {
	err := q.connect()
	if err != nil {
		// TODO 集成log
		log.Println(err)
		return
	}

	deliveries, err := q.channel.Consume(
		q.queueName, // queu name
		q.name,      // name,
		false,       // no autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		// TODO 集成log
		log.Printf("Queue Consume: %s", err)
		return
	}

	for d := range deliveries {
		go func(d amqp.Delivery) {
			defer func() {
				if err := recover(); err != nil {
					// TODO 集成log
					log.Printf("%s", err)
				}
			}()

			retry, err := consumer(context.TODO(), d.Body)
			if err != nil {
				// TODO 集成log
				log.Println(err)
			}

			if !retry {
				if err = d.Ack(true); err != nil {
					// TODO 集成log
					log.Println(err)
				}
			}
		}(d)
	}
}

func (q *RabbitMQ) Shutdown() (err error) {
	if err = q.channel.Cancel(q.name, true); err != nil {
		return errors.Wrap(err, "channel cancel failed")
	}

	if err = q.connection.Close(); err != nil {
		return errors.Wrap(err, "connection close error")
	}

	return
}

func (q *RabbitMQ) connect() (err error) {
	q.connection, err = amqp.Dial(q.url)
	if err != nil {
		return errors.Wrap(err, "amqp.Dial fail")
	}
	go func() {
		fmt.Printf("closing: %s", <-q.connection.NotifyClose(make(chan *amqp.Error)))
	}()

	q.channel, err = q.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "connection.Channel fail")
	}

	if q.opts.declareExchange {
		if err = q.channel.ExchangeDeclare(q.exchangeName, q.exchangeType, true, false, false, false, nil); err != nil {
			return errors.Wrap(err, "channel.ExchangeDeclare fail")
		}
	}

	if q.opts.declareQueue {
		if _, err = q.channel.QueueDeclare(
			q.queueName, // routing_key
			true,        // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		); err != nil {
			return errors.Wrap(err, "channel.QueueDeclare fail")
		}
	}

	if q.opts.bindQueue {
		if err = q.channel.QueueBind(q.queueName, q.routeName, q.exchangeName, false, nil); err != nil {
			return errors.Wrap(err, "channel.QueueBind fail")
		}
	}

	return
}
