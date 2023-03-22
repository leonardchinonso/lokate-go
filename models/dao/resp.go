package dao

// PlaceResp represents the TAPI response for place API
type PlaceResp struct {
	ID               string        `json:"id"`
	Member           []interface{} `json:"member"`
	View             interface{}   `json:"view"`
	Source           string        `json:"source"`
	Acknowledgements string        `json:"acknowledgements"`
}

// PublicJourneyResp represents the TAPI response for public journey API
type PublicJourneyResp struct {
	RequestTime      string        `json:"request_time"`
	Source           string        `json:"source"`
	Acknowledgements string        `json:"acknowledgements"`
	Routes           []interface{} `json:"routes"`
}
