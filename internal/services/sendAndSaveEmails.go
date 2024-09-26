package services

import (
	"EmergencyNotifictionSystem/internal/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"os"
)

func sendAndSaveEmails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.EmailRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var contacts []models.Contact
		if err := db.Where("id IN ?", request.Contacts).Find(&contacts).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		smtpHost := "smtp.gmail.com"
		smtpPort := 587
		senderEmail := os.Getenv("SENDER_EMAIL")
		senderPassword := os.Getenv("SENDER_PASSWORD")

		if senderEmail == "" || senderPassword == "" {
			log.Println("Sender email or password is not set")
			c.JSON(500, gin.H{"error": "Sender email or password is not set"})
			return
		}

		dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

		var errs []error

		for _, contact := range contacts {
			msg := gomail.NewMessage()
			msg.SetHeader("From", senderEmail)
			msg.SetHeader("To", contact.Email)
			msg.SetHeader("Subject", request.Subject)
			msg.SetBody("text/html", request.Body)

			if err := dialer.DialAndSend(msg); err != nil {
				errs = append(errs, err)
				log.Printf("Failed to send email to %s: %v", contact.Email, err)
			}
		}

		if len(errs) > 0 {
			c.JSON(500, gin.H{"error": "Some emails failed to send", "details": errs})
		} else {
			c.JSON(200, gin.H{"message": "Emails sent successfully"})
		}
	}
}
