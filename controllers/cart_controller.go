package controllers

import (
	"estore-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartController struct {
	DB *gorm.DB
}

func NewCartController(db *gorm.DB) *CartController {
	return &CartController{DB: db}
}

func (cc *CartController) AddToCart(c *gin.Context) {
	var input struct {
		ProductId uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var product models.Product
	if err := cc.DB.First(&product, input.ProductId).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	var existingCartItem models.CartItem
	if err := cc.DB.Where("customer_id = ? AND product_id = ?", customerID, input.ProductId).First(&existingCartItem).Error; err == nil {
		existingCartItem.Quantity += input.Quantity
		cc.DB.Save(&existingCartItem)

		// Preload associations before returning
		cc.DB.Preload("Product").Preload("Customer").First(&existingCartItem, existingCartItem.ID)
		c.JSON(200, gin.H{"message": "Cart updated", "data": existingCartItem})
		return
	}

	// New cart item
	cartItem := models.CartItem{
		CustomerId: customerID.(uint),
		ProductId:  input.ProductId,
		Quantity:   input.Quantity,
	}

	if err := cc.DB.Create(&cartItem).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add to cart"})
		return
	}

	var createdCartItem models.CartItem
	cc.DB.Preload("Product").Preload("Customer").First(&createdCartItem, cartItem.ID)

	c.JSON(201, gin.H{"message": "Product added to cart", "data": createdCartItem})
}

func (cc *CartController) GetCartItems(c *gin.Context) {
	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var cartItems []models.CartItem
	err := cc.DB.Preload("Product").Where("customer_id = ?", customerID).Find(&cartItems).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(200, gin.H{"data": cartItems})
}
