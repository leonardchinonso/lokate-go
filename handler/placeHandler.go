package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
)

type PlaceHandler struct {
	placeService interfaces.PlaceServiceInterface
}

// InitPlaceHandler initializes and sets up the saved places handler
func InitPlaceHandler(router *gin.Engine, version string, placeService interfaces.PlaceServiceInterface) {
	h := &PlaceHandler{
		placeService: placeService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/places")
	g := router.Group(path)

	// register endpoints
	g.POST("/", h.AddPlace)
	g.GET("/:id", h.GetPlace)
}

// AddPlace handles the request to add a place to the application
func (h *PlaceHandler) AddPlace(c *gin.Context) {
	var p dto.PlaceDTO

	// fill the contact us request by binding the JSON
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a place object
	place := dao.NewPlace(p.Type, p.Name, p.Latitude, p.Longitude, p.Accuracy, p.OSMId, p.Description)

	// add place to the database
	err := h.placeService.Create(c, place)
	if err != nil {
		log.Printf("Error creating place document in the database. Error: %v\n", err)
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusCreated("place added successfully", place)
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

// SavePlace handles the request to save a place to the application
func (h *PlaceHandler) SavePlace(c *gin.Context) {
	// get the place id from the path parameter
	_, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Failed to convert hex string to place id. Error: %v\n", err)
		resErr := errors.ErrBadRequest("invalid place id", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// retrieve the logged-in user from the authenticated request
	_, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

}
