package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type savedPlaceRepo struct {
	c *mongo.Collection
}

const savedPlaceCollectionName = "saved_places"

func NewSavedPlaceRepository(db *mongo.Database) interfaces.SavedPlaceRepositoryInterface {
	return &savedPlaceRepo{
		c: db.Collection(savedPlaceCollectionName),
	}
}

// findOneByQuery finds a document by a filter query
func (p *savedPlaceRepo) findOneByQuery(ctx context.Context, filter primitive.M, savedPlace *dao.SavedPlace) (bool, error) {
	err := p.c.FindOne(ctx, filter).Decode(savedPlace)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find saved place: %v", err)
	}
	return true, nil
}

// Create creates a new saved place document in the database
func (p *savedPlaceRepo) Create(ctx context.Context, savedPlace *dao.SavedPlace) error {
	_, err := p.c.InsertOne(ctx, savedPlace)
	if err != nil {
		return err
	}
	return nil
}

// FindByID finds a saved place by id in the database
func (p *savedPlaceRepo) FindByID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error) {
	return p.findOneByQuery(ctx, bson.M{"_id": savedPlace.Id}, savedPlace)
}

// FindOneByIDAndUserID finds a saved place by userId and the _id in the database
func (p *savedPlaceRepo) FindOneByIDAndUserID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error) {
	return p.findOneByQuery(ctx, bson.M{"user_id": savedPlace.UserId, "_id": savedPlace.Id}, savedPlace)
}

// FindOneByPlaceIDAndUserID finds a saved place by the place id and the user id
func (p *savedPlaceRepo) FindOneByPlaceIDAndUserID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error) {
	fmt.Printf("userId: %+v, placeId: %+v\n", savedPlace.UserId, savedPlace.PlaceId)
	return p.findOneByQuery(ctx, bson.M{"user_id": savedPlace.UserId, "place_id": savedPlace.PlaceId}, savedPlace)
}

// Find finds all saved places by the userId in the database
func (p *savedPlaceRepo) Find(ctx context.Context, userId primitive.ObjectID, savedPlaces *[]dao.SavedPlace) (bool, error) {
	filter := bson.M{"user_id": userId}

	cursor, err := p.c.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find saved places: %v", err)
	}

	if err = cursor.All(ctx, savedPlaces); err != nil {
		return false, fmt.Errorf("failed to find saved places: %v", err)
	}

	return true, nil
}

// updateByQuery updates a savedPlace by a specified query
func (p *savedPlaceRepo) updateByQuery(ctx context.Context, filter primitive.D, update primitive.D) error {
	_, err := p.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a savedPlace in the database
func (p *savedPlaceRepo) Update(ctx context.Context, savedPlace *dao.SavedPlace) error {
	filter := bson.D{{"_id", savedPlace.Id}, {"user_id", savedPlace.UserId}}
	update := bson.D{{"$set", bson.D{
		{"name", savedPlace.Name}, {"place_alias", savedPlace.PlaceAlias}, {"updated_at", savedPlace.UpdatedAt},
	}}}
	return p.updateByQuery(ctx, filter, update)
}

// SetAlias sets the alias of a document by a filter in the database
func (p *savedPlaceRepo) SetAlias(ctx context.Context, savedPlace *dao.SavedPlace, newAlias dao.PlaceAlias) error {
	filter := bson.D{{"user_id", savedPlace.UserId}, {"place_alias", savedPlace.PlaceAlias}}
	update := bson.D{{"$set", bson.D{
		{"place_alias", newAlias},
	}}}
	return p.updateByQuery(ctx, filter, update)
}

// Delete deletes a savedPlace by id from the database
func (p *savedPlaceRepo) Delete(ctx context.Context, savedPlace *dao.SavedPlace) error {
	filter := bson.D{{"_id", savedPlace.Id}}
	_, err := p.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
