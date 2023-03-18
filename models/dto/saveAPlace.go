package dto

import (
	"github.com/leonardchinonso/lokate-go/models/dao"
)

// SaveAPlace holds the struct details for saving a place
type SaveAPlace struct {
	UserId     string `json:"user_id"`
	PlaceId    string `json:"place_id"`
	PlaceAlias string `json:"place_alias"`
	Name       string `json:"name"`
}

// ToSavedPlaceDAO converts a SaveAPlaceDTO to SavedPlaceDAO
func ToSavedPlaceDAO(placeDTO *SaveAPlace) (*dao.SavedPlace, error) {
	// convert the dto string place alias to the dao placeAlias type
	placeAlias, err := dao.StringToPlaceAlias(placeDTO.PlaceAlias)
	if err != nil {
		return nil, err
	}

	return &dao.SavedPlace{
		Name:       placeDTO.Name,
		PlaceAlias: placeAlias,
	}, nil
}
