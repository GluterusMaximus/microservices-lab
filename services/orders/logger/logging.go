package logger

import (
	"github.com/GluterusMaximus/ci/services/orders/repository"
	"log"
	"time"
)


type LoggingService struct {
	next repository.IOrders
}

func NewLoggingService(next repository.IOrders) repository.IOrders {
	return &LoggingService{
		next: next,
	}
}

func (s *LoggingService) Get(id int) (repository.Order, error) {
	defer func(start time.Time) {
		log.Printf("getting order -> took=%v, id=%d\n", time.Since(start), id)
	}(time.Now())

	return s.next.Get(id)
}

func (s *LoggingService) Add(order repository.Order) (repository.Order, error) {
	defer func(start time.Time) {
		log.Printf("adding order -> took=%v, requestor email=%s, status=%s\n", time.Since(start), order.RequestorEmail, order.Status)
	}(time.Now())

	return s.next.Add(order)
}

func (s *LoggingService) Update(order repository.Order) (repository.Order, error) {
	defer func(start time.Time) {
		log.Printf("updating order -> took=%v, requestor email=%s, status=%s\n", time.Since(start), order.RequestorEmail, order.Status)
	}(time.Now())

	return s.next.Update(order)
}

func (s *LoggingService) Delete(order repository.Order) error {
	defer func(start time.Time) {
		log.Printf("deleting order -> took=%v, requestor email=%s, status=%s\n", time.Since(start), order.RequestorEmail, order.Status)
	}(time.Now())

	return s.next.Delete(order)
}
