package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

// SavedPlaceHandler handles saved place related requests
type SavedPlaceHandler struct {
	placeService      interfaces.PlaceServiceInterface
	savedPlaceService interfaces.SavedPlaceServiceInterface
	tokenService      interfaces.TokenServiceInterface
}

// InitSavedPlaceHandler initializes and sets up the saved places handler
func InitSavedPlaceHandler(router *gin.Engine, version string, placeService interfaces.PlaceServiceInterface, savedPlaceService interfaces.SavedPlaceServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &SavedPlaceHandler{
		placeService:      placeService,
		savedPlaceService: savedPlaceService,
		tokenService:      tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/saved-places")
	g := router.Group(path)

	// register endpoints
	g.POST("/:id", middlewares.AuthorizeUser(h.tokenService), h.AddToSavedPlaces)
	g.GET("/:id", middlewares.AuthorizeUser(h.tokenService), h.GetSavedPlace)
	g.GET("/", middlewares.AuthorizeUser(h.tokenService), h.GetSavedPlaces)
	g.PUT("/:id", middlewares.AuthorizeUser(h.tokenService), h.EditSavedPlace)
	g.DELETE("/:id", middlewares.AuthorizeUser(h.tokenService), h.DeleteSavedPlace)
}

// AddToSavedPlaces handles the request to save a place to the application
func (h *SavedPlaceHandler) AddToSavedPlaces(c *gin.Context) {
	// get the place id from the path parameter
	placeId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Failed to convert hex string to place id. Error: %v\n", err)
		resErr := errors.ErrBadRequest("invalid place id", nil)
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

	// create a dto to receive data from the request
	sap := &dto.SaveAPlace{}

	// fill the saveAPlace request by binding the JSON
	if err := c.ShouldBindJSON(sap); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), err)
		c.JSON(resErr.Status, resErr)
		return
	}

	// convert savedPlace dto to dao
	savedPlace, err := dto.ToSavedPlaceDAO(sap)
	if err != nil {
		log.Printf("Failed to convert saveAPlace DTO to savedPlaceDAO. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// add the userId and placeId to the savedPlace DAO
	savedPlace.UserId = user.Id
	savedPlace.PlaceId = placeId
	savedPlace.CreatedAt = utils.CurrentPrimitiveTime()
	savedPlace.UpdatedAt = utils.CurrentPrimitiveTime()

	// call the save place service to save the place to the application
	if err := h.savedPlaceService.AddSavedPlace(c, savedPlace); err != nil {
		log.Printf("Failed to create a place in the application")
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("place saved successfully", savedPlace)
	c.JSON(resp.Status, resp)
}

// GetSavedPlace handles the request to get a saved place
func (h *SavedPlaceHandler) GetSavedPlace(c *gin.Context) {
	// get the saved place id from the path parameter
	savedPlaceId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Failed to convert hex string to saved place id. Error: %v\n", err)
		resErr := errors.ErrBadRequest("invalid saved place id", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// retrieve the logged-in user from the authenticated request
	user, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a dao to get the data from the service
	savedPlace := &dao.SavedPlace{Id: savedPlaceId, UserId: user.Id}

	// get the saved place from the saved place service
	err = h.savedPlaceService.GetSavedPlace(c, savedPlace)
	if err != nil {
		log.Printf("Failed to retrieve saved place from savedPlaceService")
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("saved place retrieved successfully", savedPlace)
	c.JSON(resp.Status, resp)
}

// GetSavedPlaces handles the request to get all saved places
func (h *SavedPlaceHandler) GetSavedPlaces(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	user, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	// create the dao savedPlace slice to fetch the data
	var savedPlaces []dao.SavedPlace

	// get the list of saved places from the savedPlaceService
	err := h.savedPlaceService.GetSavedPlaces(c, user.Id, &savedPlaces)
	if err != nil {
		log.Printf("Failed to retrieve saved places from savedPlaceService")
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("saved places retrieved successfully", savedPlaces)
	c.JSON(resp.Status, resp)
}

// EditSavedPlace handles the request to edit a saved place
func (h *SavedPlaceHandler) EditSavedPlace(c *gin.Context) {
	// get the place id from the path parameter
	savedPlaceId, err := primitive.ObjectIDFromHex(c.Param("id"))
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
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a savedPlace object to get the data with
	savedPlace := &dao.SavedPlace{Id: savedPlaceId, UserId: user.Id}

	// fill the saveAPlace request by binding the JSON
	if err := c.ShouldBindJSON(savedPlace); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), err)
		c.JSON(resErr.Status, resErr)
		return
	}

	// call the savedPlaceService to handle it
	err = h.savedPlaceService.EditSavedPlace(c, savedPlace)
	if err != nil {
		log.Printf("Error editing a place in the savedPlaceService: %v", err)
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("saved place updated successfully", savedPlace)
	c.JSON(resp.Status, resp)
}

// DeleteSavedPlace handles the request to delete a saved place
func (h *SavedPlaceHandler) DeleteSavedPlace(c *gin.Context) {
	// get the place id from the path parameter
	savedPlaceId, err := primitive.ObjectIDFromHex(c.Param("id"))
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

	// create a savedPlace object to get the data with
	savedPlace := &dao.SavedPlace{Id: savedPlaceId, UserId: user.Id}

	// call the savedPlaceService to handle it
	err = h.savedPlaceService.DeleteSavedPlace(c, savedPlace)
	if err != nil {
		log.Printf("Error deleting a place in the savedPlaceService: %v", err)
		c.JSON(errors.Status(err), gin.H{"errors": err})
		return
	}

	resp := utils.ResponseStatusOK("deleted saved place successfully", nil)
	c.JSON(resp.Status, resp)
}
