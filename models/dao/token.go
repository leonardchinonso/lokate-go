package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leonardchinonso/lokate-go/models"
)

type TokenDAO struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id"`
	AccessToken  string             `json:"access_token" bson:"access_token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
}

// Delete deletes a token document
func (t *TokenDAO) Delete() error {
	filter := bson.D{{"user_id", t.UserId}}
	_, err := models.TokenCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

// Upsert inserts a token document if it does not exist and updates it if it does
func (t *TokenDAO) Upsert() error {
	filter := bson.D{{"user_id", t.UserId}}
	update := bson.D{{"$set", bson.D{{"user_id", t.UserId}, {"refresh_token", t.RefreshToken}, {"access_token", t.AccessToken}}}}
	opts := options.Update().SetUpsert(true)
	_, err := models.TokenCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
