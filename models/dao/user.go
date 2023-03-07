package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/leonardchinonso/lokate-go/models"
	"github.com/pkg/errors"
)

type UserDAO struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name" binding:"required" bson:"first_name"`
	LastName    string             `json:"last_name" binding:"required" bson:"last_name"`
	DisplayName string             `json:"display_name" binding:"required" bson:"display_name"`
	Email       string             `json:"email" binding:"required" bson:"email"`
	Password    string             `json:"password,omitempty" binding:"required" bson:"password"`
}

// NewUser formats the user details and creates a new user
func NewUser(firstName, lastName, email, password string) *UserDAO {
	caser := cases.Title(language.English)
	fn, ln := caser.String(firstName), caser.String(lastName)
	dn := fmt.Sprintf("%s %s", fn, ln)
	return &UserDAO{
		FirstName:   fn,
		LastName:    ln,
		DisplayName: dn,
		Email:       email,
		Password:    password,
	}
}

func (u *UserDAO) Create() error {
	_, err := models.UserCollection.InsertOne(context.TODO(), *u)
	if err != nil {
		return err
	}
	return nil
}

// FindOneByEmail returns true if a user exists and false otherwise
// returns an error if there is one
func (u *UserDAO) FindOneByEmail() (bool, error) {
	err := models.UserCollection.FindOne(context.TODO(), bson.M{"email": u.Email}).Decode(u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	return true, nil
}
