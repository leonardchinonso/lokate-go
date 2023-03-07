package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/services"
	"github.com/leonardchinonso/lokate-go/utils"
)

// Signup handles the incoming signup request and queries the Signup service
func Signup(c *gin.Context) {
	var sr dto.SignupRequest

	// fill the signup request from binding the JSON request
	if err := c.ShouldBindJSON(&sr); err != nil {
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the signup request for invalid fields
	if errs := sr.Validate(); len(errs) > 0 {
		resErr := errors.ErrBadRequest("invalid signup request", errors.ErrorToStringSlice(errs))
		c.JSON(resErr.Status, resErr)
		return
	}

	// get the user, access and refresh tokens if the user is a valid user
	user, accessToken, refreshToken, err := services.Signup(sr.FirstName, sr.LastName, sr.Email, sr.Password)
	if err != nil {
		// if the user is invalid, return an unauthorized error
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a signup (login) response and return it to the controller's caller
	loginResp := dto.NewLoginResponse(user, accessToken, refreshToken)
	resp := utils.ResponseOK("logged in successfully", loginResp)

	c.JSON(resp.Status, resp)
}

func Login(c *gin.Context) {
	var lr dto.LoginRequest

	// fill the login request from binding the JSON request
	if err := c.ShouldBindJSON(&lr); err != nil {
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the login request for invalid fields
	if errs := lr.Validate(); len(errs) > 0 {
		resErr := errors.ErrBadRequest("invalid login request", errs)
		c.JSON(resErr.Status, resErr)
		return
	}

	// get the user, access and refresh tokens if the user is a valid user
	user, accessToken, refreshToken, err := services.Login(lr.Email, lr.Password)
	if err != nil {
		// if the user is invalid, return an unauthorized error
		resErr := errors.ErrUnauthorized(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a login response and return it to the controller's caller
	loginResp := dto.NewLoginResponse(user, accessToken, refreshToken)
	resp := utils.ResponseOK("logged in successfully", loginResp)

	c.JSON(resp.Status, resp)
}

func Logout(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	u, ok := c.Get("user")
	if !ok {
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// convert the type to a user type
	user := u.(*dao.UserDAO)

	// attempt to log the user out
	err := services.Logout(user.Id)
	if err != nil {
		resErr := errors.ErrInternalServerError(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	resp := utils.ResponseOK("logged out successfully", nil)
	c.JSON(resp.Status, resp)
}
