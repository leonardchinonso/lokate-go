package handler

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
)

// PlaceHandler represents the router handler object for place requests
type PlaceHandler struct {
	placeService            interfaces.PlaceServiceInterface
	savedPlaceService       interfaces.SavedPlaceServiceInterface
	lastVisitedPlaceService interfaces.LastVisitedPlaceServiceInterface
	tapiService             interfaces.TAPIServiceInterface
	tokenService            interfaces.TokenServiceInterface
}

// InitPlaceHandler initializes and sets up the saved places handler
func InitPlaceHandler(router *gin.Engine, version string,
	placeService interfaces.PlaceServiceInterface,
	savedPlaceService interfaces.SavedPlaceServiceInterface,
	lastVisitedPlace interfaces.LastVisitedPlaceServiceInterface,
	tapiService interfaces.TAPIServiceInterface,
	tokenService interfaces.TokenServiceInterface,
) {
	h := &PlaceHandler{
		placeService:            placeService,
		savedPlaceService:       savedPlaceService,
		lastVisitedPlaceService: lastVisitedPlace,
		tapiService:             tapiService,
		tokenService:            tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/places")
	g := router.Group(path)

	// register endpoints for places
	g.POST("/", h.AddPlace)
	g.GET("/:id", h.GetPlace)

	// register endpoints for last visited places
	g.POST("/:id/last", middlewares.AuthorizeUser(h.tokenService), h.AddLastVisitedPlace)
	g.GET("/last/:num", middlewares.AuthorizeUser(h.tokenService), h.GetLastNVisitedPlaces)

	// register endpoints for search
	g.GET("/search", h.Search)
}

// AddPlace handles the request to add a place to the application
func (h *PlaceHandler) AddPlace(c *gin.Context) {
	var p dto.PlaceDTO

	// fill the add place request by binding the JSON
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a place object
	place := dao.NewPlace(p.Type, p.Name, p.Description, p.OSMId, p.ATCOCode, p.StationCode, p.TiplocCode, p.SMSCode, p.Accuracy, p.Latitude, p.Longitude)

	// make the call to the place service
	err := h.placeService.Create(c, place)
	if err != nil {
		log.Printf("Error creating a place with placeService. Error: %v\n", err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusCreated("place added successfully", nil)
	c.JSON(resp.Status, resp)
}

// GetPlace handles the request to get a place from the application
func (h *PlaceHandler) GetPlace(c *gin.Context) {
	// get the place id from the path parameter
	placeId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Failed to convert hex string to place id. Error: %v\n", err)
		resErr := errors.ErrBadRequest("invalid place id", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a place object for retrieving the place
	place := &dao.Place{Id: placeId}

	// retrieve the place from the database
	err = h.placeService.GetPlace(c, place)
	if err != nil {
		log.Printf("Error getting place document from the database. Error: %v\n", err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusOK("place retrieved successfully", place)
	c.JSON(resp.Status, resp)
}

// AddLastVisitedPlace handles the request to add a place to the last visited places
func (h *PlaceHandler) AddLastVisitedPlace(c *gin.Context) {
	// get the place id from the path parameter
	placeId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Failed to convert hex string to saved_place id. Error: %v\n", err)
		resErr := errors.ErrBadRequest("invalid saved_place id", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// retrieve the logged-in user from the authenticated request
	user, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	// create a new last visited place object to pass to the service
	lastVisitedPlace := dao.NewLastVisitedPlace(user.Id, placeId)

	// make a call to the last visited place service
	err = h.lastVisitedPlaceService.AddLastVisitedPlace(c, lastVisitedPlace)
	if err != nil {
		log.Printf("Error adding a place to last visited with placeService. Error: %v\n", err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusOK("place added successfully", nil)
	c.JSON(resp.Status, resp)
}

// GetLastNVisitedPlaces handles the request the N last visited places by a user
func (h *PlaceHandler) GetLastNVisitedPlaces(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	user, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	// read the num from the path parameter
	num := c.Param("num")

	// convert the number to an integer
	numberOfPlaces, err := strconv.Atoi(num)
	if err != nil {
		log.Printf("Failed to convert number: %v to string", num)
		resErr := errors.ErrInternalServerError("failed to get last %v visited places", num)
		c.JSON(resErr.Status, gin.H{"error": resErr.Message})
		return
	}

	// create a list of last visited places to read from the service
	var lastVisitedPlaces []dao.LastVisitedPlace

	// make a call to the place service to fetch the places
	err = h.lastVisitedPlaceService.GetLastNVisitedPlaces(c, user.Id, &lastVisitedPlaces, int64(numberOfPlaces))
	if err != nil {
		log.Printf("Error getting %v last visited places from service. Error %v", numberOfPlaces, err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusOK("last visited places retrieved successfully", lastVisitedPlaces)
	c.JSON(resp.Status, resp)
}

// Search handles the request to search for a place with text
func (h *PlaceHandler) Search(c *gin.Context) {
	// read the search value from the path parameter
	searchStr := c.Query("query")

	// create list of places object to put data into
	var places []dao.Place

	// make a call to the place service to search for the data
	placesResp, err := h.tapiService.SearchPlace(searchStr, &places)
	if err != nil {
		log.Printf("Error searching for places in the place service. Error: %v\n", err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusOK("places retrieved successfully", placesResp)
	c.JSON(resp.Status, resp)
}
