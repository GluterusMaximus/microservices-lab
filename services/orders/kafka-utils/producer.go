package kafkautils

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	globalProducer *KafkaProducer
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic: topic,
	})

	globalProducer = &KafkaProducer{
		writer: writer,
	}
}

func GetGlobalProducer() *KafkaProducer {
	return globalProducer
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
