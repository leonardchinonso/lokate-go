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
	"log"
)

// UserHandler represents the router handler object for the user requests
type UserHandler struct {
	userService  interfaces.UserServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitUserHandler initializes the user handler
func InitUserHandler(router *gin.Engine, version string, userService interfaces.UserServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/user")
	g := router.Group(path)

	g.PUT("/update-profile", middlewares.AuthorizeUser(h.tokenService), h.UpdateProfile)
}

// UpdateProfile handles the request to update user details
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	u, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	var epr dto.EditProfileRequest
	// fill the edit profile request from binding the JSON request
	if err := c.ShouldBindJSON(&epr); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the login request for invalid fields
	if errs := epr.Validate(); len(errs) > 0 {
		log.Printf("Failed to validate request. Errors: %+v", errs)
		resErr := errors.ErrBadRequest("invalid request", errs)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a new user and set their details
	user := dao.NewUser(epr.FirstName, epr.LastName, string(epr.Email), "")
	user.Id = u.Id
	user.PhoneNumber = epr.PhoneNumber

	// start the signup process
	err := h.userService.EditUserProfile(c, user)
	if err != nil {
		log.Printf("Failed to sign user up. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("profile edited successfully", user)
	c.JSON(resp.Status, resp)
}
