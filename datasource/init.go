package datasource

import (
	"log"

	"github.com/leonardchinonso/lokate-go/config"
)

// DataSource contains the DatabaseContext and the config from the environment variables
type DataSource struct {
	*DatabaseContext
	Cfg *map[string]string
}

// InitDataSource initializes the data sources
func InitDataSource() (*DataSource, error) {
	log.Println("Initializing Data Sources...")

	configMap, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	dbCtx, err := InitDB(configMap)
	if err != nil {
		return nil, err
	}

	return &DataSource{
		DatabaseContext: dbCtx,
		Cfg:             configMap,
	}, nil
}

// Close closes all the data sources and releases resources held
func (ds *DataSource) Close() {
	// cancel the context after the database client has closed its connections
	defer ds.CancelFunc()

	defer func() {
		// disconnect from the client
		if err := ds.Client.Disconnect(ds.Ctx); err != nil {
			log.Printf("Error trying to close data sources correctly: %v", err)
			panic(err)
		}
	}()
}
