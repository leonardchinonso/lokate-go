package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/models/dao"
)

// UserFromRequest gets a user set by the authentication middleware
func UserFromRequest(c *gin.Context) (*dao.User, bool) {
	u, ok := c.Get("user")
	if !ok {
		return nil, false
	}

	// convert the type to ah user type
	user := u.(*dao.User)
	return user, true
}
