package rabbitmq

import (
	"context"
	"fmt"
	srcLog "log"
	"net/http"
	"strconv"
	"time"

	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
	"gin-frame/libraries/util/dir"
	"gin-frame/libraries/util/random"
	"gin-frame/libraries/xhop"

	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
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

func GetIsLog() bool {
	cfg := config.GetConfig("log", "rabbitmq_open")

	logCfg, err := cfg.Key("turn").Bool()
	util.Must(err)
	return logCfg
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
		XHop:    nil,
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

	if GetIsLog() == true {
		writeInfoLog(logFormat, d, consumeMsg)
	}
}

func writeInfoLog(logFormat *log.LogFormat, d amqp.Delivery, consumeMsg ConsumeMsg) {
	runLogSection := "amqp"
	runLogConfig := config.GetConfig("log", runLogSection)
	runLogdir := runLogConfig.Key("dir").String()
	runLogArea, _ := runLogConfig.Key("area").Int()

	file := dir.CreateHourLogFile(runLogdir, "")
	file = file + "/" + strconv.Itoa(random.RandomN(runLogArea))

	log.Init(&log.LogConfig{
		File:           file,
		Path:           runLogdir,
		Mode:           1,
		AsyncFormatter: false,
		Debug:          true,
	}, runLogdir, file)
	log.Info(logFormat, map[string]interface{}{
		"msg": string(d.Body),
	})
}

func NewProducer(msg, amqpURI, exchangeName, exchangeType, queueName, routeName, tag string, ctx context.Context, header http.Header) error {
	var (
		parent        = opentracing.SpanFromContext(ctx)
		operationName = "producer"
		statement     = fmt.Sprintf("amqpUri:%s, exchange:%s, exchange_type:%s, queue:%s, route_name:%s, msg:%s ",
			amqpURI, exchangeName, exchangeType, queueName, routeName, msg)
		span = func() opentracing.Span {
			if parent == nil {
				return opentracing.StartSpan(operationName)
			}
			return opentracing.StartSpan(operationName, opentracing.ChildOf(parent.Context()))
		}()
		logFormat = log.LogHeaderFromContext(ctx)
		startAt   = time.Now()
		endAt     time.Time
	)
	var err error

	lastModule := logFormat.Module
	lastStartTime := logFormat.StartTime
	lastEndTime := logFormat.EndTime
	lastXHop := logFormat.XHop
	defer func() {
		logFormat.Module = lastModule
		logFormat.StartTime = lastStartTime
		logFormat.EndTime = lastEndTime
		logFormat.XHop = lastXHop
	}()

	defer span.Finish()
	defer func() {
		endAt = time.Now()

		logFormat.StartTime = startAt
		logFormat.EndTime = endAt
		latencyTime := logFormat.EndTime.Sub(logFormat.StartTime).Microseconds() // 执行时间
		logFormat.LatencyTime = latencyTime
		logFormat.XHop = xhop.NewXhopNull()

		span.SetTag("error", err != nil)
		span.SetTag("db.type", "sql")
		span.SetTag("db.statement", statement)

		if err != nil {
			log.Errorf(logFormat, "%s:[%s], error: %s", operationName, statement, err)
		} else if GetIsLog() == true {
			log.Infof(logFormat, statement)
		}

		logFormat.Module = "databus/rabbitmq"
	}()

	if parent == nil {
		span = opentracing.StartSpan("redisDo")
	} else {
		span = opentracing.StartSpan("redisDo", opentracing.ChildOf(parent.Context()))
	}
	defer span.Finish()

	span.SetTag("db.type", "redis")
	span.SetTag("db.statement", statement)
	span.SetTag("error", err != nil)

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
