package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// LastVisitedPlace is the data access object for lastVisited
type LastVisitedPlace struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	PlaceId   primitive.ObjectID `json:"place_id" bson:"place_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewLastVisitedPlace returns a new LastVisitedPlace object
func NewLastVisitedPlace(userId, placeId primitive.ObjectID) *LastVisitedPlace {
	return &LastVisitedPlace{
		UserId:    userId,
		PlaceId:   placeId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
