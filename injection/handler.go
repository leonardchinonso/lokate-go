package injection

import (
	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/lokate-go/config"
	"github.com/leonardchinonso/lokate-go/handler"
)

// injectHandlers initializes the dependencies and creates them as a config
func injectHandlers(router *gin.Engine, cfg *map[string]string, handlerCfg *HandlerConfig) {
	// get the current version number for correct routing
	version := (*cfg)[config.Version]

	// initialize the auth handler
	handler.InitAuthHandler(router, version, handlerCfg.UserService, handlerCfg.TokenService)
}
