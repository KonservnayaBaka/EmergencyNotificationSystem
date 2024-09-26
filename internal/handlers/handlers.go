package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/registration", registration(db))
	r.POST("/authorization", authorization(db))
	r.POST("/addNewContact", addNewContact(db))
	r.POST("/sendEmail", sendEmail(db))
	r.POST("/uploadContacts", uploadContacts(db))
	r.POST("/sendAndSaveEmails", sendAndSaveEmails(db))
}
