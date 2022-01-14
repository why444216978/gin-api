package kafka

import (
	"github.com/Shopify/sarama"
	"strings"
	"testing"
)

const (
	brokerAddr = "localhost:9092"
	topic      = "my_topic"
	msg        = "test_message"
)

func TestKafkaSendMessage(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	addr := strings.Split(brokerAddr, "m")
	t.Log("send message")
	SendMessage(addr, config, topic, sarama.ByteEncoder(msg))
}

func TestKafkaReceiveMessage(t *testing.T) {
	addr := strings.Split(brokerAddr, "m")
	t.Log("receive message")
	// sarama.OffsetNewest：从当前的偏移量开始消费，sarama.OffsetOldest：从最老的偏移量开始消费
	Consumer(addr, topic, 0, sarama.OffsetNewest)
}
