package models

import (
	"context"
	"fmt"
	"time"

	"github.com/leonardchinonso/lokate-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// userCollectionName is the name of the users collection in the database
	userCollectionName = "users"
	// tokenCollectionName is the name of the tokens collection in the database
	tokenCollectionName = "tokens"
)

var (
	// UserCollection is the collection of the users in the database
	UserCollection *mongo.Collection
	// TokenCollection is the collection of the users in the database
	TokenCollection *mongo.Collection
)

// ConnectDatabase takes in a URI string, connects the mongo server to it and returns the mongo client
func ConnectDatabase(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	// use the context to set the deadline for connection initiation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// fetch the client from the mongo connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		// if there is an error, cancel the context before returning to avoid a context/memory leak
		defer cancel()

		return nil, nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return client, ctx, cancel, nil
}

// Ping sends a ping request through the mongoDB client to ensure a connection has been made
func Ping(client *mongo.Client, ctx context.Context) error {
	// mongo.Client has a Ping method to ping mongoDB, the deadline of the Ping method
	// will be determined by ctx. Ping method return error if any occurred, then the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	fmt.Printf("Connected to database client successfully")

	return nil
}

func SyncCollections(client *mongo.Client) {
	// get the database connection from the client
	database := client.Database(config.Map[config.DatabaseName])

	// sync the collection from the database to the global collection holders
	UserCollection = database.Collection(userCollectionName)
	TokenCollection = database.Collection(tokenCollectionName)
}
