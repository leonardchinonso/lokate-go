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

	// initialize the handlers
	handler.InitAuthHandler(router, version, handlerCfg.UserService, handlerCfg.TokenService)
	handler.InitCommsHandler(router, version, handlerCfg.CommsService, handlerCfg.TokenService)
	handler.InitPlaceHandler(router, version, handlerCfg.PlaceService, handlerCfg.SavedPlaceService,
		handlerCfg.LastVisitedPlaceService, handlerCfg.TAPIService, handlerCfg.TokenService)
	handler.InitSavedPlaceHandler(router, version, handlerCfg.PlaceService, handlerCfg.SavedPlaceService, handlerCfg.TokenService)
	handler.InitJourneyHandler(router, version, handlerCfg.TAPIService, handlerCfg.TokenService)
	handler.InitUserHandler(router, version, handlerCfg.UserService, handlerCfg.TokenService)
}
