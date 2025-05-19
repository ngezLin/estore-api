package models

import "gorm.io/gorm"

type TransactionItem struct {
	gorm.Model
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
