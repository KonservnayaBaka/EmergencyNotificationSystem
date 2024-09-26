package main

import (
	"EmergencyNotifictionSystem/internal/handlers"
	"EmergencyNotifictionSystem/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=90814263 dbname=EmergencySystem port=1488 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Contact{})
	r := gin.Default()

	db.AutoMigrate(&models.User{})

	handlers.InitializeRoutes(r, db)

	r.Run("localhost:8000")
}
