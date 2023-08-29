package kafkautils

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	globalConsumer *KafkaConsumer
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

type MessageHandler func([]byte)

func NewKafkaConsumer(brokers []string, topic string, groupID string, handler MessageHandler) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	consumer := &KafkaConsumer{
		reader: reader,
	}

	globalConsumer = consumer

	go consumer.consumeMessages(handler)
}

func GetGlobalConsumer() *KafkaConsumer {
	return globalConsumer
}

func (c *KafkaConsumer) consumeMessages(handler MessageHandler) {
	for {
		message, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}
		handler(message.Value)
	}
}

func (c *KafkaConsumer) Close() {
	c.reader.Close()
}
