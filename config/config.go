package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Map map[string]string

const (
	// ATExpiresIn is the global config name for the AT_EXPIRES_IN variable
	ATExpiresIn = "AT_EXPIRES_IN"
	// RTExpiresIn is the global config name for the RT_EXPIRES_IN variable
	RTExpiresIn = "RT_EXPIRES_IN"
	// ATSecretKey is the global config name for the AT_SECRET_KEY variable
	ATSecretKey = "AT_SECRET_KEY"
	// RTSecretKey is the global config name for the RT_SECRET_KEY variable
	RTSecretKey = "RT_SECRET_KEY"

	// TAPIAppId is the global config name for the TAPI APP_ID variable
	TAPIAppId = "TAPI_APP_ID"
	// TAPIAppKey is the global config name for the TAPI APP_KEY variable
	TAPIAppKey = "TAPI_APP_KEY"

	// DatabaseName is the global config name for the DATABASE_NAME variable
	DatabaseName = "DATABASE_NAME"
	// BaseUri  is the global config name for the BASE_URI variable
	BaseUri = "BASE_URI"
	// Version is the global config name for the VERSION variable
	Version = "VERSION"
)

// TAPIConfig holds config variables for the transport API environment
type TAPIConfig struct {
	AppId  string
	AppKey string
}

// getEnv retrieves the value of a given key from the environment variables set
func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("failed to get value for key %v", key)
}

// InitConfig loads the config variables into the application and populates the config map with the values
func InitConfig() (*map[string]string, error) {
	// loads values from the .env file into the application
	if err := godotenv.Load("./config/dev.env"); err != nil {
		return nil, fmt.Errorf("failed to load config variables: %v", err)
	}

	Map = make(map[string]string)
	var defaultConfig = []string{TAPIAppId, TAPIAppKey, DatabaseName, ATExpiresIn, RTExpiresIn, ATSecretKey, RTSecretKey, Version, BaseUri}

	// iterate the preset config variables, retrieve their values and set them in the config map
	for _, c := range defaultConfig {
		v, err := getEnv(c)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve env variables: %v", err)
		}
		Map[c] = v
	}

	return &Map, nil
}
