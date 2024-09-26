package handlers

import (
	"EmergencyNotifictionSystem/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registration(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var existingUser models.User
		if err := db.Where("Email = ? OR login = ?", user.Email, user.Login).First(&existingUser).Error; err == nil {
			c.JSON(400, gin.H{"error": "User with this email or login already exists", "registration": false})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User registered successfully", "user": user, "registration": true})
	}
}

func authorization(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		var existingUser models.User
		if err := db.Where("login = ? AND password = ?", user.Login, user.Password).First(&existingUser).Error; err != nil {
			c.JSON(200, gin.H{
				"authorized": false,
				"Error":      err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"authorized": true,
		})
	}
}
