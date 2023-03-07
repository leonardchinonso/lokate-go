package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/lokate-go/middlewares"
)

var r *gin.Engine

func StartRouter() error {
	// start the default gin engine with the logger and middleware already loaded
	r = gin.Default()

	r.Use(middlewares.Logger())

	// map all the urls to their respective controller functions
	mapAuthUrls()

	// run the gin server on the default port :8080
	return r.Run()
}
