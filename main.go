package main

import (
	"context"
	"fmt"
	"github.com/leonardchinonso/lokate-go/config"
	"github.com/leonardchinonso/lokate-go/models"
	"github.com/leonardchinonso/lokate-go/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

const baseURI = "mongodb://localhost:27017/"

func closeAll(client *mongo.Client, ctx context.Context, cancelFunc context.CancelFunc) {
	// cancel the context after the database client has closed its connections
	defer cancelFunc()

	defer func() {
		// disconnect from the client
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func main() {
	// initialise the configuration variables
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// connect to the mongo database using the uri specified
	client, ctx, cancel, err := models.ConnectDatabase(baseURI + config.Map[config.DatabaseName])
	if err != nil {
		panic(err)
	}

	// release resource when the main function is returned
	defer closeAll(client, ctx, cancel)

	// ping the database to make sure all connections were successful
	err = models.Ping(client, ctx)
	if err != nil {
		panic(err)
	}

	// sync all the collections and models with the application
	models.SyncCollections(client)

	// initiate and run the router engine
	if err = routes.StartRouter(); err != nil {
		fmt.Println(fmt.Errorf("failed to start gin engine: %v", err))
	}
}
