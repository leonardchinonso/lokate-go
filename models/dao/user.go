package dao

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// User is the user data access object
type User struct {
	Id          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	FirstName   string              `json:"first_name" binding:"required" bson:"first_name"`
	LastName    string              `json:"last_name" binding:"required" bson:"last_name"`
	DisplayName string              `json:"display_name" binding:"required" bson:"display_name"`
	Email       string              `json:"email" binding:"required" bson:"email"`
	PhoneNumber string              `json:"phone_number" bson:"phone_number"`
	Password    string              `json:"password,omitempty" binding:"required" bson:"password"`
	CreatedAt   primitive.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.Timestamp `json:"updated_at" bson:"updated_at"`
}

// NewUser formats the user details and creates a new user
func NewUser(firstName, lastName, email, password string) *User {
	caser := cases.Title(language.English)
	fn, ln := caser.String(firstName), caser.String(lastName)
	dn := fmt.Sprintf("%s %s", fn, ln)

	currTime := utils.CurrentPrimitiveTime()

	return &User{
		FirstName:   fn,
		LastName:    ln,
		DisplayName: dn,
		Email:       email,
		Password:    password,
		CreatedAt:   currTime,
		UpdatedAt:   currTime,
	}
}
