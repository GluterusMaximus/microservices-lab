package repository

import (
	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	ID          uint `gorm:"primarykey"`
	Name        string
	Price       uint
	IsAvailable bool
}
