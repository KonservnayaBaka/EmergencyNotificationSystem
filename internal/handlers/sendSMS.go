package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func sendSMS(db *gorm.DB, number string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(400, gin.H{
			"Error": "Блокнули",
		})

	}
}
