package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GluterusMaximus/ci/services/goods/handlers"
	"github.com/GluterusMaximus/ci/services/goods/logger"
	"github.com/GluterusMaximus/ci/services/goods/repository"
	"github.com/GluterusMaximus/ci/services/goods/repository/postgres"
	utils "github.com/GluterusMaximus/ci/services/goods/utils"
	"github.com/gorilla/mux"
	ps "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------------------------------- INITIALIZING DATABASE -------------------------------
func initDB(config *utils.Config) *gorm.DB {
	dbConnector := fmt.Sprintf("postgres://%s:%s@%s/%s", config.PostgreSQL.User, config.PostgreSQL.Password, config.PostgreSQL.Host, config.PostgreSQL.Db)

	db, err := gorm.Open(ps.Open(dbConnector), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&repository.Goods{})
	return db
}

// ------------------------------- INITIALIZING KAFKA -------------------------------
func initKafka(config *utils.Config) *utils.KafkaProducer {
	return utils.NewKafkaProducer([]string{config.Kafka.Host}, config.Kafka.ProducerTopic)
}

// ------------------------------- INITIALIZING HTTP SERVER -------------------------------
func initHTTPServer(producer *utils.KafkaProducer, config *utils.Config, db *gorm.DB) {
	goods := postgres.New(db)
	goodsLogger := logger.NewLoggingService(goods)

	container := utils.NewDependencyContainer()

	// registration of goods and producer
	container.RegisterDBConnection(goods)
	container.RegisterKafkaProducer(producer)

	goodsHandler := handlers.New(goodsLogger)

	r := mux.NewRouter()

	r.HandleFunc("/api/goods", goodsHandler.GetList).Methods("GET")
	r.HandleFunc("/api/goods/{id}", goodsHandler.Get).Methods("GET")
	r.HandleFunc("/api/goods", goodsHandler.Create).Methods("POST")
	r.HandleFunc("/api/goods", goodsHandler.Update).Methods("PUT")
	r.HandleFunc("/api/goods/{id}", goodsHandler.Delete).Methods("DELETE")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPPort), r); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := utils.LoadConfigFromEnv()

	db := initDB(&config)
	producer := initKafka(&config)

	initHTTPServer(producer, &config, db)
}
