package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// SendMessage 发送消息
func SendMessage(brokerAddr []string, config *sarama.Config, topic string, value sarama.Encoder) {
	producer, err := sarama.NewSyncProducer(brokerAddr, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = producer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: value,
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("partition:%d offset:%d\n", partition, offset)
}

// Consumer 消费消息
func Consumer(brokenAddr []string, topic string, partition int32, offset int64) {
	consumer, err := sarama.NewConsumer(brokenAddr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = partitionConsumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("partition:%d offset:%d key:%s val:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
}
