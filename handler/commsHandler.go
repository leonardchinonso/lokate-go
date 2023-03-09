package handler

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/utils"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

type CommsHandler struct {
	commsService interfaces.CommsServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitCommsHandler initializes and sets up the comms handler
func InitCommsHandler(router *gin.Engine, version string, commsService interfaces.CommsServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &CommsHandler{
		commsService: commsService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/comms")
	g := router.Group(path)

	// register endpoints
	g.POST("contact-us", middlewares.AuthorizeUser(h.tokenService), h.ContactUs)
}

func (ch *CommsHandler) ContactUs(c *gin.Context) {
	// retrieve the logged-in user from the authenticated request
	user, ok := UserFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve user from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	var cu dto.ContactUsDTO

	// fill the contact us request by binding the JSON
	if err := c.ShouldBindJSON(&cu); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a contactUs object
	contactUs := dao.NewContactUs(user.Id, user.Email, cu.Subject, cu.Message)

	// send email to the app's email
	err := ch.commsService.SendContactUsEmail(c, contactUs)
	if err != nil {
		c.JSON(errors.Status(err), gin.H{"error": err})
		return
	}

	resp := utils.ResponseStatusCreated("email sent successfully", nil)
	c.JSON(resp.Status, resp)
}