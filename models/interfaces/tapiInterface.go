package interfaces

import (
	"github.com/leonardchinonso/lokate-go/models/api"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
)

// TAPIServiceInterface is the interface for the Transport API service
type TAPIServiceInterface interface {
	SearchPlace(searchStr string, places *[]dao.Place) ([]api.PlaceResponse, error)
	PublicJourneyLonLat(from dto.Location, to dto.Location) (interface{}, error)
	PublicJourneyPostcode(from dto.Postcode, to dto.Postcode) (interface{}, error)
}
