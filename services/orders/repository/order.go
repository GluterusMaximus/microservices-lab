package repository

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	RequestorEmail string
	Items     []*Item `gorm:"many2many:order_items;"`
	Status    string
}

type Item struct {
	gorm.Model
	ItemId   int
	Quantity int
	Sum      float64
}
