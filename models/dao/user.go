package dao

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// UserDAO is the user data access object
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
