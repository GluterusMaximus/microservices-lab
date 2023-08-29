package logger

import (
	"context"
	"log"
	"time"

	"github.com/GluterusMaximus/ci/services/goods/repository"
)

type LoggingService struct {
	next repository.IGoods
}

func NewLoggingService(next repository.IGoods) repository.IGoods {
	return &LoggingService{
		next: next,
	}
}

func (s *LoggingService) GetList(ctx context.Context) ([]repository.Goods, error) {
	defer func(start time.Time) {
		log.Printf("getting list -> took=%v, err=%v\n", time.Since(start), ctx.Err())
	}(time.Now())

	return s.next.GetList(ctx)
}

func (s *LoggingService) Get(id uint) (repository.Goods, error) {
	defer func(start time.Time) {
		log.Printf("getting by id -> took=%v, id=%d", time.Since(start), id)
	}(time.Now())

	return s.next.Get(id)
}

func (s *LoggingService) Add(goods repository.Goods) (repository.Goods, error) {
	defer func(start time.Time) {
		log.Printf("adding -> took=%v, id=%d, name=%s, price=%d, isAvailable=%v", time.Since(start), goods.ID, goods.Name, goods.Price, goods.IsAvailable)
	}(time.Now())

	return s.next.Add(goods)
}

func (s *LoggingService) Update(goods repository.Goods) (repository.Goods, error) {
	defer func(start time.Time) {
		log.Printf("updating -> took=%v, id=%d, name=%s, price=%d, isAvailable=%v", time.Since(start), goods.ID, goods.Name, goods.Price, goods.IsAvailable)
	}(time.Now())

	return s.next.Update(goods)
}

func (s *LoggingService) Delete(goods repository.Goods) error {
	defer func(start time.Time) {
		log.Printf("deleting -> took=%v, id=%d, name=%s, price=%d, isAvailable=%v", time.Since(start), goods.ID, goods.Name, goods.Price, goods.IsAvailable)
	}(time.Now())

	return s.next.Delete(goods)
}
