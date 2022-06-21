package rabbitmq

import (
	"context"
	"encoding/json"
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

type Config struct {
	ServiceName  string
	Host         string
	Port         int
	Virtual      string
	User         string
	Pass         string
	ExchangeType string
	ExchangeName string
	QueueName    string
	RouteName    string
}

type Option struct {
	declareExchange bool
	declareQueue    bool
	bindQueue       bool
	qos             int
}

type OptionFunc func(*Option)

func defaultOption() *Option {
	return &Option{
		declareExchange: false,
		declareQueue:    false,
		bindQueue:       false,
		qos:             10,
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

func WithQos(qos int) OptionFunc {
	return func(o *Option) { o.qos = qos }
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

func New(cfg *Config, opts ...OptionFunc) (*RabbitMQ, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}

	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}

	return &RabbitMQ{
		opts:         opt,
		name:         cfg.ServiceName,
		url:          fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Virtual),
		exchangeName: cfg.ExchangeName,
		exchangeType: cfg.ExchangeType,
		queueName:    cfg.QueueName,
		routeName:    cfg.RouteName,
	}, nil
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
		true,           // set true, when no queue match Basic.Return
		false,          // set false, not dependent consumers
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

	select {
	case r := <-q.connection.NotifyClose(make(chan *amqp.Error)):
		b, _ := json.Marshal(r)
		// TODO log
		log.Println("q.connection.NotifyClose:" + string(b))
	default:
	}

	q.channel, err = q.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "connection.Channel fail")
	}

	if err = q.channel.Qos(q.opts.qos, 0, false); err != nil {
		return errors.Wrap(err, "channel.Qos fail")
	}

	select {
	case r := <-q.channel.NotifyClose(make(chan *amqp.Error)):
		b, _ := json.Marshal(r)
		// TODO log
		log.Println("q.channel.NotifyClose:" + string(b))
	case r := <-q.channel.NotifyCancel(make(chan string)):
		// TODO log
		log.Println("q.channel.NotifyCancel:" + r)
	case r := <-q.channel.NotifyReturn(make(chan amqp.Return)):
		b, _ := json.Marshal(r)
		// TODO log
		log.Println("q.channel.NotifyReturn:" + string(b))
	default:
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
