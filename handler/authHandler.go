package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	userService  interfaces.UserServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitAuthHandler initializes and sets up the auth handler
func InitAuthHandler(router *gin.Engine, version string, userService interfaces.UserServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/auth")
	g := router.Group(path)

	// register endpoints
	g.POST("/signup", h.Signup)
	g.POST("/login", h.Login)
	g.POST("/logout", middlewares.AuthorizeUser(h.tokenService), h.Logout)
}

// Signup handles the incoming signup request
func (a *AuthHandler) Signup(c *gin.Context) {
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

	// create a new user object with the details
	user := dao.NewUser(sr.FirstName, sr.LastName, string(sr.Email), string(sr.Password))

	// start the signup process
	userId, err := a.userService.Signup(c, user, sr.Password)
	if err != nil {
		log.Printf("Failed to sign user up. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	log.Printf("ID OF USER: %+v\n", userId)

	// retrieve the new user object from the database
	user, err = a.userService.GetUserByID(c, userId)
	if err != nil {
		log.Printf("Failed to get user from database. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	// create the access and refresh token pairs
	at, rt, err := a.tokenService.GenerateTokenPair(c, user)
	if err != nil {
		log.Printf("Failed to generate user token pair. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	// create a signup (login) response and return it to the handler's caller
	loginResp := dto.NewLoginResponse(*user, at, rt)
	resp := utils.ResponseStatusCreated("signed up successfully", loginResp)

	c.JSON(resp.Status, resp)
}

// Login handles the incoming login request
func (a *AuthHandler) Login(c *gin.Context) {
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

	// create a new user object with the details
	user := dao.NewUser("", "", string(lr.Email), string(lr.Password))

	// start the login process
	err := a.userService.Login(c, user, lr.Password)
	if err != nil {
		log.Printf("Failed to login user. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	// create the access and refresh token pairs
	at, rt, err := a.tokenService.GenerateTokenPair(c, user)
	if err != nil {
		log.Printf("Failed to generate user token pair. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	// create a login response and return it to the handler's caller
	loginResp := dto.NewLoginResponse(*user, at, rt)
	resp := utils.ResponseStatusCreated("logged in successfully", loginResp)

	c.JSON(resp.Status, resp)
}

// Logout handles the incoming logout request
func (a *AuthHandler) Logout(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	u, ok := c.Get("user")
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	// convert the type to a user type
	user := u.(*dao.UserDAO)
	
	// attempt to log the user out
	err := a.userService.Logout(c, user.Id)
	if err != nil {
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusOK("logged out successfully", nil)
	c.JSON(resp.Status, resp)
}
