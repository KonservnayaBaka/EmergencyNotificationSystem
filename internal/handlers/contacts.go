package handlers

import (
	"EmergencyNotifictionSystem/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addNewContact(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var contact models.Contact
		if err := c.ShouldBindJSON(&contact); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var existingContact models.Contact
		if err := db.Where("phone = ? OR email = ?", contact.Phone, contact.Email).First(&existingContact).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&contact).Error; err != nil {
					c.JSON(400, gin.H{"Error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"addNewContact": true})
			} else {
				c.JSON(500, gin.H{"Error": err.Error()})
			}
		} else {
			c.JSON(400, gin.H{"error": "Contact with this phone or email already exists", "addNewContact": false})
		}
	}
}
