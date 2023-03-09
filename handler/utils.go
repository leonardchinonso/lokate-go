package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/models/dao"
)

func UserFromRequest(c *gin.Context) (*dao.UserDAO, bool) {
	u, ok := c.Get("user")
	if !ok {
		return nil, false
	}

	// convert the type to ah user type
	user := u.(*dao.UserDAO)
	return user, true
}
