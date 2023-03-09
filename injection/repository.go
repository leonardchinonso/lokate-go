package injection

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/repository"
)

// ServicesConfig is the custom type for starting up services
type ServicesConfig struct {
	UserRepo      interfaces.UserRepositoryInterface
	TokenRepo     interfaces.TokenRepositoryInterface
	ContactUsRepo interfaces.ContactUsRepositoryInterface
}

// injectRepositories initializes the dependencies and creates them as a config for services injection
func injectRepositories(db *mongo.Database) *ServicesConfig {
	return &ServicesConfig{
		UserRepo:      repository.NewUserRepository(db),
		TokenRepo:     repository.NewTokenRepository(db),
		ContactUsRepo: repository.NewContactUsRepository(db),
	}
}
