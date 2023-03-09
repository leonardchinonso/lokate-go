package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

// ContactUsDAO is the struct for contact us data
type ContactUsDAO struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserEmail string             `json:"user_email" bson:"user_email"`
	Subject   string             `json:"subject" bson:"subject"`
	Message   string             `json:"message" bson:"message"`
}

func NewContactUs(userId primitive.ObjectID, userEmail, subject, message string) *ContactUsDAO {
	return &ContactUsDAO{
		UserId:    userId,
		UserEmail: userEmail,
		Subject:   subject,
		Message:   message,
	}
}
