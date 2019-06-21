package models

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	ProductID uint `json:"-"`
	Product   Product
	Qty       int
	Status    string
	UserID    uint
	User      User
}
