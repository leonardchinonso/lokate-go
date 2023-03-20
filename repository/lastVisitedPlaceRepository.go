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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type lastVisitedPlaceRepo struct {
	c *mongo.Collection
}

const lastVisitedPlaceCollectionName = "last_visited"

func NewLastVisitedPlaceRepository(db *mongo.Database) interfaces.LastVisitedPlaceRepositoryInterface {
	return &lastVisitedPlaceRepo{
		c: db.Collection(lastVisitedPlaceCollectionName),
	}
}

// Create creates a new document in the database
func (l *lastVisitedPlaceRepo) Create(ctx context.Context, lastVisitedPlace *dao.LastVisitedPlace) error {
	_, err := l.c.InsertOne(ctx, lastVisitedPlace)
	if err != nil {
		return err
	}
	return nil
}

// FindLastNVisitedPlaces finds the last visited N places by a user
func (l *lastVisitedPlaceRepo) FindLastNVisitedPlaces(ctx context.Context, userId primitive.ObjectID, lastVisitedPlaces *[]dao.LastVisitedPlace, N int64) (bool, error) {
	// create a filter object
	filter := bson.D{{"user_id", userId}}

	// sort the documents by creation date
	sort := bson.D{{"created_at", -1}}

	// modify the find query with the options
	opts := options.Find().SetSort(sort).SetLimit(N)

	cursor, err := l.c.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find last visited %v places: %v", N, err)
	}

	if err = cursor.All(ctx, lastVisitedPlaces); err != nil {
		return false, fmt.Errorf("failed to find last visited %v places: %v", N, err)
	}

	return true, nil
}
