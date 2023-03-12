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

const savedPlaceCollectionName = "savedPlaces"

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

// FindOneByUserID finds a saved place by userId in the database
func (p *savedPlaceRepo) FindOneByUserID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error) {
	return p.findOneByQuery(ctx, bson.M{"user_id": savedPlace.UserId}, savedPlace)
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

// Update updates a savedPlace in the database
func (p *savedPlaceRepo) Update(ctx context.Context, savedPlace *dao.SavedPlace) error {
	filter := bson.D{{"_id", savedPlace.Id}, {"user_id", savedPlace.UserId}}
	update := bson.D{{"$set", bson.D{
		{"name", savedPlace.Name}, {"use_as", savedPlace.UseAs}, {"updated_at", savedPlace.UpdatedAt},
	}}}
	_, err := p.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
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
