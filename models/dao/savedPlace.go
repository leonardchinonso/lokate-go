package dao

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UsedAs is an enum that represents what a place can be used as
type UsedAs string

const (
	Home UsedAs = "HOME"
	Work UsedAs = "WORK"
	None UsedAs = "NONE"
)

// SavedPlace is the savedPlace data access object
type SavedPlace struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	PlaceId   primitive.ObjectID  `json:"place_id" bson:"place_id,omitempty"`
	UserId    primitive.ObjectID  `json:"user_id" bson:"user_id"`
	Name      string              `json:"name" bson:"name"`
	UseAs     UsedAs              `json:"use_as" bson:"use_as"`
	CreatedAt primitive.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.Timestamp `json:"updated_at" bson:"updated_at"`
}

// Validate checks that a saved place is a valid one
func (sp *SavedPlace) Validate() error {
	switch sp.UseAs {
	case Home, Work, None:
		return nil
	}
	return fmt.Errorf("invalid useAs type for place")
}
