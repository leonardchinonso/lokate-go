package injection

import (
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/service"
)

// HandlerConfig holds the configuration values for initializing the handlers
type HandlerConfig struct {
	UserService  interfaces.UserServiceInterface
	TokenService interfaces.TokenServiceInterface
	CommsService interfaces.CommsServiceInterface
	ReqService   interfaces.RequestServiceInterface
}

// injectServices initializes the dependencies and creates them as a config for handler injection
func injectServices(cfg *map[string]string, servCfg *ServicesConfig) (*HandlerConfig, error) {
	// initialize the user service with the needed config
	userService := service.NewUserService(servCfg.UserRepo, servCfg.TokenRepo)

	// initialize the token service with the needed config
	tokenService, err := service.NewTokenService(cfg, servCfg.TokenRepo)
	if err != nil {
		return nil, err
	}

	// initialize the comms service with  the needed config
	commsService := service.NewCommsService(cfg, servCfg.ContactUsRepo)

	// initialize the external requests service with the needed config
	reqService := service.NewRequestService()

	return &HandlerConfig{
		UserService:  userService,
		TokenService: tokenService,
		CommsService: commsService,
		ReqService:   reqService,
	}, nil
}
