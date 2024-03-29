package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

type userRepo struct {
	c *mongo.Collection
}

const userCollectionName = "users"

// NewUserRepository returns a token interface with all the model repository methods
func NewUserRepository(db *mongo.Database) interfaces.UserRepositoryInterface {
	return &userRepo{
		c: db.Collection(userCollectionName),
	}
}

// Create creates a new user document in the database
func (ur *userRepo) Create(ctx context.Context, user *dao.User) (primitive.ObjectID, error) {
	result, err := ur.c.InsertOne(ctx, user)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// FindByID finds a user by id in the database
func (ur *userRepo) FindByID(ctx context.Context, user *dao.User) (bool, error) {
	err := ur.c.FindOne(ctx, bson.M{"_id": user.Id}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	return true, nil
}

// FindByEmail finds a user by email in the database
func (ur *userRepo) FindByEmail(ctx context.Context, user *dao.User) (bool, error) {
	err := ur.c.FindOne(ctx, bson.M{"email": user.Email}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	return true, nil
}

// Update updates a user in the database
func (ur *userRepo) Update(ctx context.Context, user *dao.User) error {
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{
		{"first_name", user.FirstName}, {"last_name", user.LastName},
		{"display_name", user.DisplayName}, {"email", user.Email},
		{"phone_number", user.PhoneNumber}, {"updated_at", user.UpdatedAt},
	}}}
	return ur.updateByQuery(ctx, filter, update)
}

// updateByQuery updates a savedPlace by a specified query
func (ur *userRepo) updateByQuery(ctx context.Context, filter primitive.D, update primitive.D) error {
	_, err := ur.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
