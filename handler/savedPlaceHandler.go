package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

// SavedPlaceHandler handles saved place related requests
type SavedPlaceHandler struct {
	placeService interfaces.PlaceServiceInterface
	savedPlaceService interfaces.SavedPlaceServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitSavedPlaceHandler initializes and sets up the saved places handler
func InitSavedPlaceHandler(router *gin.Engine, version string, placeService interfaces.PlaceServiceInterface, savedPlaceService interfaces.SavedPlaceServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &SavedPlaceHandler{
		placeService: placeService,
		savedPlaceService:   savedPlaceService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/saved-places")
	g := router.Group(path)

	// register endpoints
	g.POST("/", middlewares.AuthorizeUser(h.tokenService), h.AddSavedPlace)
}

// AddSavedPlace adds a new place to the user saved places
func (h *SavedPlaceHandler) AddSavedPlace(c *gin.Context) {
//	var
}
