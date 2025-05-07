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

func (ac *AuthController) AdminLogin(c *gin.Context) {
	var loginReq models.LoginRequestAdmin

	//buat check valiasi masukan dari admin
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//buat check email admin
	var admin models.Admin
	if err := ac.DB.Where("email = ?", loginReq.Email).First(&admin).Error; err != nil {
		c.JSON(401, gin.H{"Error": "Invalid Email"})
		return
	}

	if err := admin.CheckPasswordAdmin(loginReq.Password); err != nil {
		c.JSON(401, gin.H{"Error": "Invalid Password"})
		return
	}

	token, err := utils.GenerateTokenAdmin(admin.ID)
	if err != nil {
		c.JSON(500, gin.H{"Eror": "Error Generating Token"})
		return
	}
	//keluarkan output
	c.JSON(200, gin.H{
		"message": "Login Successfully!",
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

func (ac *AuthController) CustomerLogin(c *gin.Context) {
	var loginReq models.LoginRequestCustomer

	//buat check valiasi masukan dari customer
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//buat check email customer
	var customer models.Customer
	if err := ac.DB.Where("email = ?", loginReq.Email).First(&customer).Error; err != nil {
		c.JSON(401, gin.H{"Error": "Invalid Email"})
		return
	}

	if err := customer.CheckPasswordCustomer(loginReq.Password); err != nil {
		c.JSON(401, gin.H{"Error": "Invalid Password"})
		return
	}

	token, err := utils.GenerateTokenCustomer(customer.ID)
	if err != nil {
		c.JSON(500, gin.H{"Eror": "Error Generating Token"})
		return
	}
	//keluarkan output
	c.JSON(200, gin.H{
		"message": "Login Successfully!",
		"token":   token,
	})
}
