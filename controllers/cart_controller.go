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

		cc.DB.Preload("Product").Preload("Customer").First(&existingCartItem, existingCartItem.ID)
		c.JSON(200, gin.H{"message": "Cart updated", "data": existingCartItem})
		return
	}

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

func (cc *CartController) ClearCart(c *gin.Context) {
	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := cc.DB.Where("customer_id = ?", customerID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(200, gin.H{"message": "Cart cleared successfully"})
}

func (cc *CartController) UpdateCartItemQuantity(c *gin.Context) {
	productID := c.Param("product_id")
	var input struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Quantity < 1 {
		c.JSON(400, gin.H{"error": "Invalid quantity"})
		return
	}

	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var cartItem models.CartItem
	if err := cc.DB.Where("customer_id = ? AND product_id = ?", customerID, productID).First(&cartItem).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cart item not found"})
		return
	}

	cartItem.Quantity = input.Quantity
	cc.DB.Save(&cartItem)

	cc.DB.Preload("Product").Preload("Customer").First(&cartItem, cartItem.ID)

	c.JSON(200, gin.H{"message": "Cart item updated", "data": cartItem})
}

func (cc *CartController) RemoveCartItem(c *gin.Context) {
	productID := c.Param("product_id")
	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := cc.DB.Where("customer_id = ? AND product_id = ?", customerID, productID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to remove cart item"})
		return
	}

	c.JSON(200, gin.H{"message": "Cart item removed"})
}

func (cc *CartController) Checkout(c *gin.Context) {
	customerID, exists := c.Get("customerId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var cartItems []models.CartItem
	cc.DB.Preload("Product").Where("customer_id = ?", customerID).Find(&cartItems)

	if len(cartItems) == 0 {
		c.JSON(400, gin.H{"error": "Cart is empty"})
		return
	}

	var totalPrice float64
	var transactionItems []models.TransactionItem

	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.Product.Price
		transactionItems = append(transactionItems, models.TransactionItem{
			ProductID: item.Product.ID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		})

		cc.DB.Model(&item.Product).Update("stock", item.Product.Stock-item.Quantity)
	}

	transaction := models.Transaction{
		CustomerID: customerID.(uint),
		TotalPrice: totalPrice,
		Items:      transactionItems,
	}

	cc.DB.Create(&transaction)
	cc.DB.Where("customer_id = ?", customerID).Delete(&models.CartItem{})
	cc.DB.Preload("Items.Product").Preload("Customer").First(&transaction, transaction.ID)

	c.JSON(200, gin.H{
		"message": "Checkout successful",
		"data":    transaction,
	})
}
