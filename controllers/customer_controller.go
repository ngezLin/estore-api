package controllers

import (
	"estore-api/models"

	"gorm.io/gorm"
)

type CustomerController struct {
	DB *gorm.DB
}

// bikin variable buat menampung
var CustomersInMemory = []models.Admin{}

func NewCustomerController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}
