package dto

import "fmt"

// Location represent a longitude and latitude map point
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Notation  string `json:"notation"`
}

// NewLocation returns a new Location object
func NewLocation(latitude, longitude string) Location {
	return Location{
		Latitude:  latitude,
		Longitude: longitude,
		Notation:  fmt.Sprintf("%s,%s", longitude, latitude),
	}
}
