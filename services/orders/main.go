package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GluterusMaximus/ci/services/orders/handlers"
	kafkautils "github.com/GluterusMaximus/ci/services/orders/kafka-utils"
	"github.com/GluterusMaximus/ci/services/orders/logger"
	"github.com/GluterusMaximus/ci/services/orders/repository"
	"github.com/GluterusMaximus/ci/services/orders/repository/postgres"
	"github.com/GluterusMaximus/ci/services/orders/utils"
	"github.com/gorilla/mux"
	ps "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var (
// 	httpPort      int
// 	pgHost        string
// 	pgUser        string
// 	pgPass        string
// 	pgDb          string
// 	kafkaHost     string
// 	producerTopic string
// 	consumerTopic string
// )

// func init() {
// 	httpPort = 8080
// 	pgUser = os.Getenv("POSTGRES_USER")
// 	pgPass = os.Getenv("POSTGRES_PASSWORD")
// 	pgHost = os.Getenv("POSTGRES_HOST")
// 	pgDb = os.Getenv("POSTGRES_DB")
// 	kafkaHost = os.Getenv("KAFKA_HOST")
// 	producerTopic = os.Getenv("PRODUCER_TOPIC")
// 	consumerTopic = os.Getenv("CONSUMER_TOPIC")
// }

func initDB(config *utils.Config) *gorm.DB {
	dbConnector := fmt.Sprintf("postgres://%s:%s@%s/%s", config.PostgreSQL.User, config.PostgreSQL.Password, config.PostgreSQL.Host, config.PostgreSQL.Db)

	db, err := gorm.Open(ps.Open(dbConnector), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&repository.Order{})
	return db
}

func initHTTPServer(server *handlers.Server, config *utils.Config, db *gorm.DB) {
	r := mux.NewRouter()

	r.HandleFunc("/api/orders/{id}", server.Get).Methods("GET")
	r.HandleFunc("/api/orders", server.Create).Methods("POST")
	r.HandleFunc("/api/orders", server.Update).Methods("PUT")
	r.HandleFunc("/api/orders/{id}", server.Delete).Methods("DELETE")

	kafkautils.NewKafkaConsumer([]string{config.Kafka.Host}, config.Kafka.ConsumerTopic, "orders-group", func(message []byte) {
		handlers.OutOfStockHandler(message, db)
	})
	
	defer kafkautils.GetGlobalConsumer().Close()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPPort), r); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := utils.LoadConfigFromEnv()
	db := initDB(&config)

	kafkautils.NewKafkaProducer([]string{config.Kafka.Host}, config.Kafka.ProducerTopic)
	defer kafkautils.GetGlobalProducer().Close()

	

	orders := postgres.New(db)
	orderLogger := logger.NewLoggingService(orders)

	server := handlers.New(orderLogger)
	
	initHTTPServer(server, &config, db)
}