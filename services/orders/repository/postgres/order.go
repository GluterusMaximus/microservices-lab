package postgres

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/GluterusMaximus/ci/services/orders/repository"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func GetGood(id int) (map[string]interface{}, error) {
    var results map[string]interface{}

    resp, err := http.Get("http://local-goods/api/goods/" + strconv.FormatUint(uint64(id), 10))
    if err != nil {
        return nil, fmt.Errorf("connection error: %s", err.Error())
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("parse body error: %s", err.Error())
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("goods response: %s", body)
    }

    json.Unmarshal(body, &results)
    return results, nil
}

func (r *Repository) Get(id int) (repository.Order, error) {
    var order repository.Order

    result := r.db.Preload("Items").First(&order, id)
    if result.Error != nil {
        return repository.Order{}, result.Error
    }

    return order, nil
}

func (r *Repository) Add(order repository.Order) (repository.Order, error) {
    var itemsLen int = len(order.Items)

    for i := 0; i < itemsLen; i++ {
        item, err := GetGood(order.Items[i].ItemId)
        if err != nil || item == nil {
            return repository.Order{}, err
        }

        if item["IsAvailable"] == false {
            return repository.Order{}, fmt.Errorf("error creating order: good %d is not available", uint(item["ID"].(float64)))
        }

        order.Items[i].Sum = item["Price"].(float64) * float64(order.Items[i].Quantity)
        order.Status = "Processing"
    }

    r.db.Create(&order)
    return order, nil
}

func (r *Repository) Update(order repository.Order) (repository.Order, error) {
    var existingOrder repository.Order
    res := r.db.First(&existingOrder, "Id = ?", order.ID)
    if res.RowsAffected == 0 {
        return repository.Order{}, fmt.Errorf("order is not present")
    }

    r.db.Save(&order)
    return order, nil
}

func (r *Repository) Delete(order repository.Order) error {
    var existingOrder repository.Order
    res := r.db.First(&existingOrder, "Id = ?", order.ID)
    if res.RowsAffected == 0 {
        return fmt.Errorf("order is not present")
    }

    r.db.Delete(&order)
    return nil
}