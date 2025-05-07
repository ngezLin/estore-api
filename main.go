package main

import (
	"estore-api/config"
	"estore-api/models"
	"estore-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env")
	}

	r := gin.Default()
	db := config.ConnectDatabase()

	db.AutoMigrate(
		&models.Admin{},
		&models.Customer{},
	)

	routes.SetupRoutes(r, db)
	r.Run(":8080")
}
