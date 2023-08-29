package utils

import "os"

type Config struct {
	HTTPPort   int
	PostgreSQL PostgresConfig
	Kafka      KafkaConfig
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Db       string
}

type KafkaConfig struct {
	Host          string
	ProducerTopic string
	ConsumerTopic string
}

func LoadConfigFromEnv() Config {
	return Config{
		HTTPPort: 8080,
		PostgreSQL: PostgresConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Db:       os.Getenv("POSTGRES_DB"),
		},
		Kafka: KafkaConfig{
			Host:          os.Getenv("KAFKA_HOST"),
			ProducerTopic: os.Getenv("PRODUCER_TOPIC"),
			ConsumerTopic: os.Getenv("CONSUMER_TOPIC"),
		},
	}
}
