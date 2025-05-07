package controllers

import (
	"estore-api/models"

	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

// bikin variable buat menampung
var adminsInMemory = []models.Admin{}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}
