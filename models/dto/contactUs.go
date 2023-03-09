package dto

// ContactUsDTO represents the request information for contact us route
type ContactUsDTO struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}
