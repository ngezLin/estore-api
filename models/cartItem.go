package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CustomerId uint `json:"customer_id"`
	ProductId  uint `json:"product_id"`
	Quantity   int  `json:"quantity"`

	Customer Customer `gorm:"foreignKey:CustomerId" json:"customer"`
	Product  Product  `gorm:"foreignKey:ProductId" json:"product"`
}
