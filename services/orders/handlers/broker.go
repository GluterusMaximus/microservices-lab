package handlers

import (
	"encoding/json"
	"log"
	"strconv"

	kafkautils "github.com/GluterusMaximus/ci/services/orders/kafka-utils"
	"github.com/GluterusMaximus/ci/services/orders/repository"
	"gorm.io/gorm"
)

func OutOfStockHandler(message []byte, db *gorm.DB) error {
	id, err := strconv.Atoi(string(message))
	if err != nil {
		log.Printf("error parsing broker message: %v", err)
		return err
	}

	log.Printf("order id: %d\n", id)
	var orders []repository.Order
	if err := db.Distinct("orders.*").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN items ON order_items.item_id = items.id").
		Where("items.item_id IN (?)", id).
		Find(&orders).Error; err != nil {
		log.Printf("error updating orders with these goods: %v", err)
		return err
	}

	for i := 0; i < len(orders); i++ {
		orders[i].Status = "Under Review"
		db.Save(orders[i])
	}

	if err := db.Where("item_id = ?", id).Delete(&repository.Item{}).Error; err != nil {
		log.Printf("error deleting items: %v", err)
		return err
	}

	for _, order := range orders {
		messageJSON, err := json.Marshal(order)
		if err != nil {
			return err
		}

		err = kafkautils.GetGlobalProducer().Produce(string(messageJSON))
		if err != nil {
			return err
		}
	}

	return nil
}
