package service

import (
	"encoding/json"
	"fmt"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/api"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"log"

	"github.com/leonardchinonso/lokate-go/config"
	"github.com/leonardchinonso/lokate-go/datasource"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

// tapiService holds the structure for services associated with TAPI
type tapiService struct {
	tapiAppId         string
	tapiAppKey        string
	tapiPlacesUrl     string
	tapiPublicJourney string
	tapiServiceName   string
}

// NewTAPIService returns an interface for the TAPI service methods
func NewTAPIService(cfg *map[string]string) interfaces.TAPIServiceInterface {
	return &tapiService{
		tapiAppId:         (*cfg)[config.TAPIAppId],
		tapiAppKey:        (*cfg)[config.TAPIAppKey],
		tapiPlacesUrl:     (*cfg)[config.TAPIPlacesUrl],
		tapiPublicJourney: (*cfg)[config.TAPIPublicJourneyUrl],
		tapiServiceName:   (*cfg)[config.TAPIServiceName],
	}
}

// SearchPlace makes a http request to TAPI to get a query string
func (ts *tapiService) SearchPlace(searchStr string, places *[]dao.Place) ([]api.PlaceResponse, error) {
	if searchStr == "" {
		log.Printf("Failed to get data for url. Error: invalid search query\n")
		return nil, errors.ErrBadRequest("invalid search query", nil)
	}

	// build the url to search
	url := fmt.Sprintf("%s?query=%s&app_id=%s&app_key=%s", ts.tapiPlacesUrl, searchStr, ts.tapiAppId, ts.tapiAppKey)

	// create object to marshal into
	var placeResp dao.PlaceResp

	// make a http request to the url
	err := datasource.Get(url, &placeResp)
	if err != nil {
		log.Printf("Failed to get data for url: %v. Error: %v\n", url, err)
		return nil, errors.ErrInternalServerError("failed to reach TAPI", nil)
	}

	// Convert the "members" field to have access to all its nested values
	for _, item := range placeResp.Member {
		marshalled, err := json.Marshal(item)
		if err != nil {
			log.Printf("failed to marshal with error: %v\n", err)
			return nil, errors.ErrInternalServerError("failed to convert place", nil)
		}

		place := &dao.Place{}
		err = json.Unmarshal(marshalled, place)
		if err != nil {
			log.Printf("failed to unmarshal with error: %v\n", err)
			return nil, errors.ErrInternalServerError("failed to convert place", nil)
		}

		*places = append(*places, *place)
	}

	return api.ToPlacesResponses(places), nil
}

// PublicJourneyLonLat specifies the method for getting a public journey by lonlat format from TAPI
func (ts *tapiService) PublicJourneyLonLat(from dto.Location, to dto.Location) (interface{}, error) {
	// build the url to query location by lonlat format
	url := fmt.Sprintf("%s/from/lonlat:%s/to/lonlat:%s.json?service=%s&app_id=%s&app_key=%s",
		ts.tapiPublicJourney, from.Notation, to.Notation, ts.tapiServiceName, ts.tapiAppId, ts.tapiAppKey)

	return ts.publicJourney(url)
}

// PublicJourneyPostcode specifies the method for getting a public journey by postcode from TAPI
func (ts *tapiService) PublicJourneyPostcode(from dto.Postcode, to dto.Postcode) (interface{}, error) {
	// build the url to query location by postcode format
	url := fmt.Sprintf("%s/from/postcode:%s/to/postcode:%s.json?service=%s&app_id=%s&app_key=%s",
		ts.tapiPublicJourney, from.StrVal, to.StrVal, ts.tapiServiceName, ts.tapiAppId, ts.tapiAppKey)

	return ts.publicJourney(url)
}

// publicJourney gets a public journey from TAPI using a specified url
func (ts *tapiService) publicJourney(url string) (interface{}, error) {
	// create object to marshal into
	var pubJourneyResp dao.PublicJourneyResp

	// make a http request to the url
	err := datasource.Get(url, &pubJourneyResp)
	if err != nil {
		log.Printf("Failed to get data for url: %v. Error: %v", url, err)
		return nil, errors.ErrInternalServerError("failed to reach TAPI", nil)
	}

	return pubJourneyResp, nil
}
