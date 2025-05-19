package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	CustomerID uint `json:"customer_id"`
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	TotalPrice float64 `json:"total_price"`
	Items []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}