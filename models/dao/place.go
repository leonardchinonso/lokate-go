package dao

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Place is the place data access object
type Place struct {
	Id          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Type        string              `json:"type" bson:"type"`
	Name        string              `json:"name" bson:"name"`
	Latitude    string              `json:"latitude" bson:"latitude"`
	Longitude   string              `json:"longitude" bson:"longitude"`
	Accuracy    string              `json:"accuracy" bson:"accuracy"`
	OSMId       string              `json:"osm_id" bson:"osm_id"`
	Description string              `json:"description" bson:"description"`
	Key         string              `json:"key" bson:"key"`
	CreatedAt   primitive.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.Timestamp `json:"updated_at" bson:"updated_at"`
}

func NewPlace(placeType, name, lat, lon, acc, osmid, desc string) *Place {
	return &Place{
		Type:        placeType,
		Name:        name,
		Latitude:    lat,
		Longitude:   lon,
		Accuracy:    acc,
		OSMId:       osmid,
		Description: desc,
		Key:         fmt.Sprintf("%s%s", lat, lon),
		CreatedAt:   utils.CurrentPrimitiveTime(),
		UpdatedAt:   utils.CurrentPrimitiveTime(),
	}
}
