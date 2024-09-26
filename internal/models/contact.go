package models

type Contact struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type EmailRequest struct {
	Contacts []uint `json:"contacts"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}
