package routes

import (
	"estore-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	// adminController := controllers.NewAdminController(db)

	//public
	r.POST("/admin/register", authController.AdminRegister)
	r.POST("/customer/register", authController.CustomerRegister)
}
