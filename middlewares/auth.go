package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

// AuthorizeUser reads the token from the request header and gets the logged-in user
func AuthorizeUser(ts interfaces.TokenServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenArr := c.Request.Header["Token"]
		if len(tokenArr) == 0 || tokenArr[0] == "" {
			resErr := errors.ErrUnauthorized("empty token value", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		splitTokenStr := strings.Split(tokenArr[0], "Bearer ")
		if len(splitTokenStr) != 2 {
			resErr := errors.ErrUnauthorized("must provide Authorization header with format `Bearer {token}`", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		user, err := ts.UserFromAccessToken(splitTokenStr[1])
		if err != nil {
			resErr := errors.ErrUnauthorized("sorry, you're not authorized for this request", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
