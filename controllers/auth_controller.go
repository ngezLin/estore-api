package controllers

import (
	"estore-api/models"
	"estore-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) AdminRegister(c *gin.Context) {
	var admin models.Admin

	//validasi input
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//hash password
	if err := admin.HashPassword(admin.Password); err != nil {
		c.JSON(500, gin.H{"Error": err.Error()})
		return
	}

	//insert ke database
	result := ac.DB.Create(&admin)
	if result.Error != nil {
		c.JSON(400, gin.H{"Error": "Error creating admin"})
		return
	}

	//generate token jwt
	token, err := utils.GenerateTokenAdmin(admin.ID)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Error generating token"})
		return
	}

	//output
	c.JSON(201, gin.H{
		"message": "Admin created successfully",
		"token":   token,
	})
}

func (ac *AuthController) CustomerRegister(c *gin.Context) {
	var customer models.Customer

	//validasi input
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//hash password
	if err := customer.HashPassword(customer.Password); err != nil {
		c.JSON(500, gin.H{"Error": err.Error()})
		return
	}

	//insert ke database
	result := ac.DB.Create(&customer)
	if result.Error != nil {
		c.JSON(400, gin.H{"Error": "Error creating customer"})
		return
	}

	//generate token jwt
	token, err := utils.GenerateTokenCustomer(customer.ID)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Error generating token"})
		return
	}

	//output
	c.JSON(201, gin.H{
		"message": "Customer created successfully",
		"token":   token,
	})
}
