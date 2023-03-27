package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/middlewares"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
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
	path := fmt.Sprintf("%s%s", version, "/journey")
	g := router.Group(path)

	g.PATCH("/change-profile", middlewares.AuthorizeUser(h.tokenService), h.ChangeProfile)
}

// ChangeProfile handles the request to change a user details
func (h *UserHandler) ChangeProfile(c *gin.Context) {

}
