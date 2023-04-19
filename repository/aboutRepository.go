package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

// aboutRepo holds the collection for the about mongo collection
type aboutRepo struct {
	c *mongo.Collection
}

const aboutCollectionName = "about"

// NewAboutRepository returns a contactUs interface with all the model repository methods
func NewAboutRepository(db *mongo.Database) interfaces.AboutRepositoryInterface {
	return &aboutRepo{
		c: db.Collection(aboutCollectionName),
	}
}

// GetDetails retrieves the about details from the database
func (a *aboutRepo) GetDetails(ctx context.Context, details *dao.About) error {
	return a.findOneByQuery(ctx, bson.M{}, details)
}

// findOneByQuery finds a document by a filter query
func (a *aboutRepo) findOneByQuery(ctx context.Context, filter primitive.M, about *dao.About) error {
	err := a.c.FindOne(ctx, filter).Decode(about)
	if err != nil {
		return fmt.Errorf("failed to get details from database: %v", err)
	}
	return nil
}
