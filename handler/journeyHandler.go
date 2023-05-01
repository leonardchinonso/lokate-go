package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
	"log"
)

// JourneyHandler represents the router handler object for journey requests
type JourneyHandler struct {
	tapiService  interfaces.TAPIServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitJourneyHandler initializes the journey handler
func InitJourneyHandler(router *gin.Engine, version string, tapiService interfaces.TAPIServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &JourneyHandler{
		tapiService:  tapiService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/journey")
	g := router.Group(path)

	g.GET("/lonlat", h.PublicJourneyLonLat)
	g.GET("/postcode", h.PublicJourneyPostcode)
}

// PublicJourneyLonLat handles the request to get journeys using lonlat format
func (h *JourneyHandler) PublicJourneyLonLat(c *gin.Context) {
	// read the from and to values from the path parameter
	fromLat := c.Query("from_lat")
	fromLon := c.Query("from_lon")
	toLat := c.Query("to_lat")
	toLon := c.Query("to_lon")

	// convert the values to a location type
	fromLoc := dto.NewLocation(fromLat, fromLon)
	toLoc := dto.NewLocation(toLat, toLon)

	// get the journey from the TAPI service
	journeyResp, err := h.tapiService.PublicJourneyLonLat(fromLoc, toLoc)
	if err != nil {
		log.Printf("Error getting public journey for latitude and longitude values. Error: %v", err)
		c.JSON(errors.Status(err), gin.H{"errors": err})
		return
	}

	resp := utils.ResponseStatusOK("routes retrieved successfully", journeyResp)
	c.JSON(resp.Status, resp)
}

// PublicJourneyPostcode handles the request to get journeys using postcode format
func (h *JourneyHandler) PublicJourneyPostcode(c *gin.Context) {
	// read the from and to values from the path parameter
	fromPostcode := c.Query("from")
	toPostcode := c.Query("to")

	// convert the values to a postcode type
	from, err := dto.NewPostcode(fromPostcode)
	if err != nil {
		log.Printf("Error getting public journey for postcode values. Error: %v", err)
		resErr := errors.ErrBadRequest("invalid postcode", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr.Message})
		return
	}

	to, err := dto.NewPostcode(toPostcode)
	if err != nil {
		log.Printf("Error getting public journey for postcode values. Error: %v", err)
		resErr := errors.ErrBadRequest("invalid postcode", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr.Message})
		return
	}

	// get the journey from the TAPI service
	journeyResp, err := h.tapiService.PublicJourneyPostcode(from, to)
	if err != nil {
		log.Printf("Error getting public journey for postcode values. Error: %v", err)
		c.JSON(errors.Status(err), gin.H{"errors": err})
		return
	}

	resp := utils.ResponseStatusOK("routes retrieved successfully", journeyResp)
	c.JSON(resp.Status, resp)
}
