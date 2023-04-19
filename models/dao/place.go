package dao

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Place is the place data access object
type Place struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Latitude    *float64           `json:"latitude" bson:"latitude"`
	Longitude   *float64           `json:"longitude" bson:"longitude"`
	Accuracy    *int               `json:"accuracy,omitempty" bson:"accuracy"`
	Description string             `json:"description,omitempty" bson:"description"`
	OSMId       string             `json:"osm_id,omitempty" bson:"osm_id"`
	ATCOCode    string             `json:"atcocode,omitempty" bson:"atcocode"`
	StationCode string             `json:"station_code,omitempty" bson:"station_code"`
	TiplocCode  string             `json:"tiploc_code,omitempty" bson:"tiploc_code"`
	SMSCode     string             `json:"smscode,omitempty" bson:"smscode"`
	Distance    *int               `json:"distance,omitempty" bson:"distance"`
	Key         string             `json:"key" bson:"key"`
}

func NewPlace(placeType, name, desc, osmid, atcocode, stationCode, tiplocCode, smsCode string, acc *int, lat, lon *float64) *Place {
	return &Place{
		Type:        placeType,
		Name:        name,
		Latitude:    lat,
		Longitude:   lon,
		Accuracy:    acc,
		Description: desc,
		OSMId:       osmid,
		ATCOCode:    atcocode,
		StationCode: stationCode,
		TiplocCode:  tiplocCode,
		SMSCode:     smsCode,
		Key:         fmt.Sprintf("%f%f", *lat, *lon),
	}
}

func (p *Place) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	p.Name = m["name"].(string)

	// check type existence
	if m["type"] != nil {
		p.Type = m["type"].(string)
	}

	// check accuracy existence
	if m["accuracy"] != nil {
		p.Accuracy = m["accuracy"].(*int)
	}

	// check description existence
	if m["description"] != nil {
		p.Description = m["description"].(string)
	}

	// check osm_id existence
	if m["osm_id"] != nil {
		p.OSMId = m["osm_id"].(string)
	}

	// check atcocode existence
	if m["atcocode"] != nil {
		p.ATCOCode = m["atcocode"].(string)
	}

	// check station_code existence
	if m["station_code"] != nil {
		p.StationCode = m["station_code"].(string)
	}

	// check tiploc_code existence
	if m["tiploc_code"] != nil {
		p.TiplocCode = m["tiploc_code"].(string)
	}

	// check smscode existence
	if m["smscode"] != nil {
		p.SMSCode = m["smscode"].(string)
	}

	// check distance existence
	if m["distance"] != nil {
		p.Distance = m["distance"].(*int)
	}

	lat := m["latitude"].(float64)
	p.Latitude = &lat

	lon := m["longitude"].(float64)
	p.Longitude = &lon

	return nil
}

func (p *Place) CopyValues(place Place) {
	p.Id = place.Id
	p.Name = place.Type
	p.Name = place.Name
	p.Latitude = place.Latitude
	p.Longitude = place.Longitude
	p.Accuracy = place.Accuracy
	p.Description = place.Description
	p.OSMId = place.OSMId
	p.ATCOCode = place.ATCOCode
	p.StationCode = place.StationCode
	p.TiplocCode = place.TiplocCode
	p.SMSCode = place.SMSCode
	p.Distance = place.Distance
	p.Key = place.Key
}
