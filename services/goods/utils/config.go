package utils

import (
	"os"
)

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
}

// TODO: delete temp variables, write immediately to the Config struct
func LoadConfigFromEnv() Config {
	httpPort := 8080 // default value

	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgDb := os.Getenv("POSTGRES_DB")

	kafkaHost := os.Getenv("KAFKA_HOST")
	producerTopic := os.Getenv("PRODUCER_TOPIC")

	return Config{
		HTTPPort: httpPort,
		PostgreSQL: PostgresConfig{
			User:     pgUser,
			Password: pgPass,
			Host:     pgHost,
			Db:       pgDb,
		},
		Kafka: KafkaConfig{
			Host:          kafkaHost,
			ProducerTopic: producerTopic,
		},
	}
}
