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

type placeRepo struct {
	c *mongo.Collection
}

const placeCollectionName = "places"

func NewPlaceRepository(db *mongo.Database) interfaces.PlaceRepositoryInterface {
	return &placeRepo{
		c: db.Collection(placeCollectionName),
	}
}

// findByQuery finds a document by a filter query
func (p *placeRepo) findByQuery(ctx context.Context, filter primitive.M, place *dao.Place) (bool, error) {
	err := p.c.FindOne(ctx, filter).Decode(place)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find place: %w", err)
	}
	return true, nil
}

// Create creates a new place document in the database
func (p *placeRepo) Create(ctx context.Context, place *dao.Place) error {
	_, err := p.c.InsertOne(ctx, place)
	if err != nil {
		return err
	}
	return nil
}

// FindByID finds a place by id in the database
func (p *placeRepo) FindByID(ctx context.Context, place *dao.Place) (bool, error) {
	return p.findByQuery(ctx, bson.M{"_id": place.Id}, place)
}

// FindByKey finds a place by key in the database
func (p *placeRepo) FindByKey(ctx context.Context, place *dao.Place) (bool, error) {
	return p.findByQuery(ctx, bson.M{"key": place.Key}, place)
}

// PopulatePlacesInLastVisited populates the places array in the lastVisitedPlaces object
func (p *placeRepo) PopulatePlacesInLastVisited(ctx context.Context, lastVisitedPlaces *[]dao.LastVisitedPlace) (bool, error) {
	for idx, lastVisited := range *lastVisitedPlaces {
		var pl dao.Place
		exists, err := p.findByQuery(ctx, bson.M{"_id": lastVisited.PlaceId}, &pl)
		lastVisited.Place = pl
		(*lastVisitedPlaces)[idx] = lastVisited
		if err != nil {
			return exists, err // TODO: should throw an error if one of the places is missing
		}
	}
	return true, nil
}

// PopulatePlacesInSavedPlaces populates the places array in the savedPlaces object
func (p *placeRepo) PopulatePlacesInSavedPlaces(ctx context.Context, savedPlaces *[]dao.SavedPlace) (bool, error) {
	for idx, savedPlace := range *savedPlaces {
		var pl dao.Place
		exists, err := p.findByQuery(ctx, bson.M{"_id": savedPlace.PlaceId}, &pl)
		savedPlace.Place = pl
		(*savedPlaces)[idx] = savedPlace
		if err != nil {
			return exists, err // TODO: should throw an error if one of the places is missing
		}
	}
	return true, nil
}
