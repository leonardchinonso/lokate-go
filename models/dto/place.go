package dto

type PlaceDTO struct {
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
