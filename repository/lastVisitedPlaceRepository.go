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

// FindLastNVisitedPlacesWithAggregate finds the last visited N places by a user and aggregates the places field
func (l *lastVisitedPlaceRepo) FindLastNVisitedPlacesWithAggregate(ctx context.Context, userId primitive.ObjectID, lastVisitedPlaces *[]dao.LastVisitedPlace, N int64) (bool, error) {
	fmt.Println(userId)
	var ps []dao.Place

	// create a filter object
	filter := bson.M{"user_id": userId}

	// create a sort object
	sort := bson.M{"created_at": -1}

	var aggSearch, aggSort, aggLimit, aggPopulate, aggProject bson.M

	// filter by search term
	aggSearch = bson.M{"$match": filter}

	// sort the documents by the sort criteria
	aggSort = bson.M{"$sort": sort}

	// set limit to the first N elements
	aggLimit = bson.M{"$limit": N}

	// populate the Places field
	aggPopulate = bson.M{"$lookup": bson.M{
		"from":         "places",   // from the collection name to populate
		"localField":   "place_id", // the field to populate on the struct
		"foreignField": "_id",      // the field on the place struct to populate by
		"as":           "places",   // the field on the destination struct to populate into
	}}

	// take the first element from the populated array, array is of size 1
	aggProject = bson.M{"$project": bson.M{
		"place": bson.M{"arrayElemAt": []interface{}{"place", 0}},
	}}

	cursor, err := l.c.Aggregate(ctx, []bson.M{
		aggSearch, aggSort, aggLimit, aggPopulate, aggProject,
	})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find last visited %v places: %v", N, err)
	}

	if err = cursor.All(ctx, &ps); err != nil {
		return false, fmt.Errorf("failed to find last visited %v places: %v", N, err)
	}

	fmt.Printf("Values: %+v\n", ps)

	return true, nil
}
