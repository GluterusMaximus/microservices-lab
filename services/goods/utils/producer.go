package utils

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &KafkaProducer{
		writer: writer,
	}
}

func (p *KafkaProducer) Produce(message string) error {
	err := p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)

	if err != nil {
		log.Printf("failed to produce message: %v", err)
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() {
	p.writer.Close()
}
