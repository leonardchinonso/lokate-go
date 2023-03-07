package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Logger is a middleware function that logs details of an incoming request
// referenced from https://www.tabnine.com/blog/golang-gin/
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client IP address
		clientIP := c.ClientIP()

		// Get the current time
		now := time.Now()

		// Log the request
		log.Printf("[%s] %s %s %s", now.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, clientIP)

		// Proceed to the next handler
		c.Next()
	}
}
