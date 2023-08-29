package utils

import (
	"github.com/GluterusMaximus/ci/services/goods/repository"
)

// common di container
var (
	di *DependencyContainer
)

type DependencyContainer struct {
	dbConnection  repository.IGoods
	kafkaProducer *KafkaProducer
}

func ResolveInitedDi() *DependencyContainer {
	return di
}

func NewDependencyContainer() *DependencyContainer {
	di = &DependencyContainer{}
	return di
}

func (c *DependencyContainer) RegisterDBConnection(db repository.IGoods) {
	c.dbConnection = db
}

func (c *DependencyContainer) ResolveGoodsRepository() repository.IGoods {
	return c.dbConnection
}

func (c *DependencyContainer) RegisterKafkaProducer(producer *KafkaProducer) {
	c.kafkaProducer = producer
}

func (c *DependencyContainer) ResolveKafkaProducer() *KafkaProducer {
	return c.kafkaProducer
}
