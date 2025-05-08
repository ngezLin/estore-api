package routes

import (
	"estore-api/controllers"
	"estore-api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	productController := controllers.NewProductController(db)
	cartController := controllers.NewCartController(db)

	//public
	r.POST("/admin/register", authController.AdminRegister)
	r.POST("/customer/register", authController.CustomerRegister)
	r.POST("/admin/login", authController.AdminLogin)
	r.POST("/customer/login", authController.CustomerLogin)

	authAdmin := r.Group("/admin", middlewares.AuthMiddlewareAdmin())
	{
		authAdmin.POST("/products", productController.CreateProduct)
		authAdmin.GET("/products", productController.GetAllProducts)
		authAdmin.PUT("/products/:id", productController.UpdateProduct)
		authAdmin.DELETE("/products/:id", productController.DeleteProduct)

	}

	authCustomer := r.Group("/customer", middlewares.AuthMiddlewareCustomer())
	{
		authCustomer.GET("/products", productController.GetAvailableProducts)
		authCustomer.GET("/products/:id", productController.GetProductDetails)
		authCustomer.POST("/cart", cartController.AddToCart)
		authCustomer.GET("/cart", cartController.GetCartItems)
		authCustomer.DELETE("/cart", cartController.ClearCart)
		authCustomer.PUT("/cart/:product_id", cartController.UpdateCartItemQuantity)
	}
}
