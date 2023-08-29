package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	kafkautils "github.com/GluterusMaximus/ci/services/orders/kafka-utils"
	"github.com/GluterusMaximus/ci/services/orders/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	db repository.IOrders
}

func New(db repository.IOrders) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("failed parse order id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	order, err := s.db.Get(orderId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get order: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, jsonError := json.Marshal(order)
	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}
	w.Write(jsonResponse)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var reqOrder repository.Order

	err := json.NewDecoder(r.Body).Decode(&reqOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := s.db.Add(reqOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, jsonError := json.Marshal(order)
	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}
	w.Write(jsonResponse)
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	var reqOrder repository.Order

	err := json.NewDecoder(r.Body).Decode(&reqOrder)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse order: %s", err.Error()), http.StatusBadRequest)
		return
	}

	order, err := s.db.Update(reqOrder)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update order: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	messageJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("failed to marshal order: %v\n", err)
	}

	err = kafkautils.GetGlobalProducer().Produce(string(messageJSON))
	if err != nil {
		log.Printf("failed to produce order: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, jsonError := json.Marshal(order)
	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}
	w.Write(jsonResponse)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("failed parse order id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = s.db.Delete(repository.Order{
		Model: gorm.Model{
			ID: uint(orderId),
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete order: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order %d deleted", orderId)
}
