package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlaceDTO struct {
	Type        string              `json:"type"`
	Name        string              `json:"name"`
	Latitude    string              `json:"latitude"`
	Longitude   string              `json:"longitude"`
	Accuracy    string              `json:"accuracy"`
	OSMId       string              `json:"osm_id"`
	Description string              `json:"description"`
	CreatedAt   primitive.Timestamp `json:"created_at"`
	UpdatedAt   primitive.Timestamp `json:"updated_at"`
}
