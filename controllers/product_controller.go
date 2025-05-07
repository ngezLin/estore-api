package controllers

import (
	"estore-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

// admin
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := pc.DB.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(201, gin.H{"message": "Product created successfully", "data": product})
}

func (pc *ProductController) GetAllProducts(c *gin.Context) {
	var products []models.Product
	pc.DB.Find(&products)
	c.JSON(200, gin.H{"data": products})
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pc.DB.Save(&product)
	c.JSON(200, gin.H{"message": "Product updated", "data": product})
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := pc.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted"})
}

// customer
func (pc *ProductController) GetAvailableProducts(c *gin.Context) {
	var products []models.Product
	if err := pc.DB.Where("Stock > ?", 0).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieved products"})
		return
	}
	c.JSON(200, gin.H{"data": products})
}

func (pc *ProductController) GetProductDetails(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(200, gin.H{"data": product})
}
