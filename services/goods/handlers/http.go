package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GluterusMaximus/ci/services/goods/repository"
	utils "github.com/GluterusMaximus/ci/services/goods/utils"
	"github.com/gorilla/mux"
)

type Server struct {
	db repository.IGoods
}

func New(db repository.IGoods) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) GetList(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	goods, err := s.db.GetList(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get goods: %s", err.Error()), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, jsonError := json.Marshal(goods)

	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}

	w.Write(jsonResponse)
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	goodId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("failed parse good id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	good, err := s.db.Get(uint(goodId))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get good: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, jsonError := json.Marshal(good)

	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}

	w.Write(jsonResponse)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var reqGood repository.Goods

	err := json.NewDecoder(r.Body).Decode(&reqGood)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse good: %s", err.Error()), http.StatusBadRequest)
		return
	}

	good, err := s.db.Add(reqGood)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create good: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, jsonError := json.Marshal(good)

	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}

	w.Write(jsonResponse)
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	var reqGood repository.Goods

	err := json.NewDecoder(r.Body).Decode(&reqGood)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse good: %s", err.Error()), http.StatusBadRequest)
		return
	}

	good, err := s.db.Update(reqGood)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update good: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, jsonError := json.Marshal(good)

	if jsonError != nil {
		http.Error(w, fmt.Sprintf("failed to send response: %s", err.Error()), http.StatusBadRequest)
	}

	w.Write(jsonResponse)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	goodId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 4)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed parse good id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = s.db.Delete(repository.Goods{
		ID: uint(goodId)})

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to good order: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	go utils.ResolveInitedDi().ResolveKafkaProducer().Produce(mux.Vars(r)["id"])
	fmt.Fprintf(w, "Good %d deleted", goodId)
}
