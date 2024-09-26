package handlers

import (
	"EmergencyNotifictionSystem/internal/models"
	"encoding/csv"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
)

func uploadContacts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var contacts []models.Contact

		switch file.Header.Get("Content-Type") {
		case "text/csv":
			f, err := file.Open()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			contacts, err = parseCSV(f)
		case "application/xml":
			f, err := file.Open()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			contacts, err = parseXML(f)
		default:
			c.JSON(400, gin.H{"error": "Unsupported file type"})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&contacts).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Contacts uploaded successfully"})
	}
}

func parseCSV(file io.Reader) ([]models.Contact, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var contacts []models.Contact
	for _, record := range records {
		contacts = append(contacts, models.Contact{
			Name:  record[0],
			Phone: record[1],
			Email: record[2],
		})
	}

	return contacts, nil
}

func parseXML(file io.Reader) ([]models.Contact, error) {
	var contacts []models.Contact
	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

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
