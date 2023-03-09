package repository

import (
	"context"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
)

type commsRepo struct {
	c *mongo.Collection
}

const contactUsCollectionName = "contact_us"

func NewContactUsRepository(db *mongo.Database) interfaces.ContactUsRepositoryInterface {
	return &commsRepo{
		c: db.Collection(contactUsCollectionName),
	}
}

func (c *commsRepo) Create(ctx context.Context, comms *dao.ContactUsDAO) error {
	_, err := c.c.InsertOne(ctx, comms)
	if err != nil {
		return err
	}
	return nil
}
