package dao

import (
	"github.com/leonardchinonso/lokate-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContactUs is the struct for contact us data
type ContactUs struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID  `json:"user_id" bson:"user_id"`
	UserEmail string              `json:"user_email" bson:"user_email"`
	Subject   string              `json:"subject" bson:"subject"`
	Message   string              `json:"message" bson:"message"`
	CreatedAt primitive.Timestamp `json:"created_at" bson:"created_at"`
}

func NewContactUs(userId primitive.ObjectID, userEmail, subject, message string) *ContactUs {
	return &ContactUs{
		UserId:    userId,
		UserEmail: userEmail,
		Subject:   subject,
		Message:   message,
		CreatedAt: utils.CurrentPrimitiveTime(),
	}
}
