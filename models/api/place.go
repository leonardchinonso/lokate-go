package api

import "github.com/leonardchinonso/lokate-go/models/dao"

// PlaceResponse is a struct for the API Response of a place
type PlaceResponse struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Accuracy    *int     `json:"accuracy"`
	Description string   `json:"description"`
	OSMId       string   `json:"osm_id"`
	ATCOCode    string   `json:"atcocode"`
	StationCode string   `json:"station_code"`
	TiplocCode  string   `json:"tiploc_code"`
	SMSCode     string   `json:"smscode"`
	Distance    *int     `json:"distance"`
}

// NewPlaceResponse returns a new place response
func NewPlaceResponse(p *dao.Place) *PlaceResponse {
	return &PlaceResponse{
		Type:        p.Type,
		Name:        p.Name,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		Accuracy:    p.Accuracy,
		Description: p.Description,
		OSMId:       p.OSMId,
		ATCOCode:    p.ATCOCode,
		StationCode: p.StationCode,
		TiplocCode:  p.TiplocCode,
		SMSCode:     p.SMSCode,
		Distance:    p.Distance,
	}
}

// ToPlacesResponses converts a data access place to the api response type
func ToPlacesResponses(places *[]dao.Place) []PlaceResponse {
	var placeResponses = make([]PlaceResponse, len(*places))

	for i, p := range *places {
		placeResponse := NewPlaceResponse(&p)
		placeResponses[i] = *placeResponse
	}

	return placeResponses
}
