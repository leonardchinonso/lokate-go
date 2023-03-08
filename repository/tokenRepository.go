package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

type tokenRepo struct {
	c *mongo.Collection
}

const tokenCollectionName = "tokens"

// NewTokenRepository returns a token interface with all the model repository methods
func NewTokenRepository(db *mongo.Database) interfaces.TokenRepositoryInterface {
	return &tokenRepo{
		c: db.Collection(tokenCollectionName),
	}
}

// Upsert updates a token by the userId if it exists in the database
// it inserts a new document if it does not exist
func (tr *tokenRepo) Upsert(ctx context.Context, token *dao.TokenDAO) error {
	filter := bson.D{{"user_id", token.UserId}}
	update := bson.D{{"$set", bson.D{{"user_id", token.UserId}, {"refresh_token", token.RefreshToken}, {"access_token", token.AccessToken}}}}
	opts := options.Update().SetUpsert(true)
	_, err := tr.c.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a token from the token collection
func (tr *tokenRepo) Delete(ctx context.Context, userId primitive.ObjectID) error {
	filter := bson.D{{"user_id", userId}}
	_, err := tr.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
