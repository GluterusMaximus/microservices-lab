// Package kafkautils provides utility functions for working with Kafka consumers.
package kafkautils

import (
	"context"

	config "github.com/GluterusMaximus/ci/services/notify/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	logger *logrus.Logger
}

type MessageHandler func([]byte)

func NewKafkaConsumer(brokers []string, topic string, groupID string, apiKey string, logger *logrus.Logger, handler MessageHandler) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: config.MinKafkaBytes,
		MaxBytes: config.MaxKafkaBytes,
	})

	consumer := &KafkaConsumer{
		reader: reader,
		logger: logger,
	}

	go consumer.consumeMessages(handler)

	return consumer
}

func (c *KafkaConsumer) consumeMessages(handler MessageHandler) {
	for {
		message, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			c.logger.Errorf("failed to read message: %v", err)
			continue
		}
		c.logger.Infof("Received message: %s", message.Value)

		handler(message.Value)
	}
}

func (c *KafkaConsumer) Close() {
	c.reader.Close()
}
