package repository

import "context"

type IGoods interface {
	GetList(ctx context.Context) ([]Goods, error)
	Get(id uint) (Goods, error)
	Add(Goods) (Goods, error)
	Update(Goods) (Goods, error)
	Delete(Goods) error
}
