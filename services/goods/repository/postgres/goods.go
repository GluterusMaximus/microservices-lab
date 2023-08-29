package postgres

import (
	"context"
	"fmt"

	"github.com/GluterusMaximus/ci/services/goods/repository"
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

func (r *Repository) GetList(ctx context.Context) ([]repository.Goods, error) {
	var goods []repository.Goods
	result := r.db.Find(&goods)

	if result.Error != nil {
		return nil, result.Error
	}

	return goods, nil
}

func (r *Repository) Get(id uint) (repository.Goods, error) {
	var good repository.Goods
	result := r.db.First(&good, id)

	if result.Error != nil {
		return repository.Goods{}, result.Error
	}

	return good, nil
}

func (r *Repository) Add(goods repository.Goods) (repository.Goods, error) {
	if goods.ID != 0 {
		return repository.Goods{}, fmt.Errorf("goods already exist")
	}

	r.db.Create(&goods)
	return goods, nil
}

func (r *Repository) Update(goods repository.Goods) (repository.Goods, error) {
	var existingGoods repository.Goods
	result := r.db.First(&existingGoods, "Id = ?", goods.ID)

	if result.RowsAffected == 0 {
		return repository.Goods{}, fmt.Errorf("goods is not present")
	}

	r.db.Save(&goods)
	return goods, nil
}

func (r *Repository) Delete(goods repository.Goods) error {
	var existingGoods repository.Goods
	result := r.db.First(&existingGoods, "Id = ?", goods.ID)

	if result.RowsAffected == 0 {
		return fmt.Errorf("goods is not present")
	}

	r.db.Delete(&goods)
	return nil
}
