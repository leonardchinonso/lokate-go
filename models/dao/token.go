package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token is the token data access object
type Token struct {
	Id           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserId       primitive.ObjectID  `json:"user_id" bson:"user_id"`
	AccessToken  string              `json:"access_token" bson:"access_token"`
	RefreshToken string              `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    primitive.Timestamp `json:"created_at" bson:"created_at"`
}
