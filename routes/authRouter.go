package routes

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/controllers"
	"github.com/leonardchinonso/lokate-go/middlewares"
)

func mapAuthUrls() {
	const routePrefix = "auth"

	r.POST(fmt.Sprintf("%s/signup", routePrefix), controllers.Signup)
	r.POST(fmt.Sprintf("%s/login", routePrefix), controllers.Login)

	r.POST(fmt.Sprintf("%s/logout", routePrefix), middlewares.Authenticate, controllers.Logout)
}
