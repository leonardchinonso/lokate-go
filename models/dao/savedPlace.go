package dao

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type PlaceAlias string

const (
	Home PlaceAlias = "HOME"
	Work PlaceAlias = "WORK"
	None PlaceAlias = "NONE"
)

// StringToPlaceAlias converts a string to PlaceAlias type after validation
func StringToPlaceAlias(s string) (PlaceAlias, error) {
	alias := PlaceAlias(strings.ToUpper(s))
	err := alias.Validate()
	if err != nil {
		return "", err
	}
	return alias, nil
}

// Validate makes sure a PlaceAlias type is valid
func (pa PlaceAlias) Validate() error {
	switch pa {
	case Home, Work, None:
		return nil
	}

	return fmt.Errorf("invalid place alias: %v", pa)
}

// IsHome determines if a placeAlias has the home value
func (pa PlaceAlias) IsHome() bool {
	return pa == Home
}

// IsWork determines if a placeAlias has the work value
func (pa PlaceAlias) IsWork() bool {
	return pa == Work
}

// IsNone determines if a placeAlias has the none value
func (pa PlaceAlias) IsNone() bool {
	return pa == None
}

// SavedPlace is the savedPlace data access object
type SavedPlace struct {
	Id         primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	PlaceId    primitive.ObjectID  `json:"place_id" bson:"place_id,omitempty"`
	Place      Place               `json:"place" bson:"place,omitempty"`
	UserId     primitive.ObjectID  `json:"user_id" bson:"user_id"`
	Name       string              `json:"name" bson:"name"`
	PlaceAlias PlaceAlias          `json:"place_alias" bson:"place_alias"`
	CreatedAt  primitive.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt  primitive.Timestamp `json:"updated_at" bson:"updated_at"`
}
