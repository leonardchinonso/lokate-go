package services

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
)

// Signup handles the user creation and logs the user in
func Signup(firstName, lastName string, email dto.Email, password dto.Password) (dao.UserDAO, string, string, error) {
	// hash the password to hide its real value
	hashedPassword, err := password.Hash()
	if err != nil {
		return dao.UserDAO{}, "", "", fmt.Errorf("failed to sign up the user: %v", err)
	}

	// create a new user object
	user := dao.NewUser(firstName, lastName, string(email), hashedPassword)

	// check that the email is not taken
	userExists, err := user.FindOneByEmail()
	if err != nil {
		return dao.UserDAO{}, "", "", fmt.Errorf("failed to fetch user details: %v", err)
	}

	// if the user exists, return an error saying so
	if userExists {
		fmt.Printf("User IS: %+v\n", user)
		return dao.UserDAO{}, "", "", fmt.Errorf(errors.ErrEmailTaken)
	}

	// create a new user with the credentials
	err = user.Create()
	if err != nil {
		return dao.UserDAO{}, "", "", fmt.Errorf("failed to create the user: %v", err)
	}

	return Login(email, password)
}

// Login logs the user into the application and returns the authentication tokens
func Login(email dto.Email, password dto.Password) (dao.UserDAO, string, string, error) {
	// find the user by email and password
	user := &dao.UserDAO{Email: string(email)}
	userExists, err := user.FindOneByEmail()
	if err != nil { // if an unexpected error occurs
		return dao.UserDAO{}, "", "", fmt.Errorf("failed to fetch user details: %v", err)
	}

	// if the user does not exist, then the password and/or email are wrong
	if !userExists {
		fmt.Printf("User Does Not Exist\n")
		return dao.UserDAO{}, "", "", fmt.Errorf(errors.ErrInvalidLogin)
	}

	// if the user exists, but the password is not correct
	if !password.IsEqualHash(user.Password) {
		return dao.UserDAO{}, "", "", fmt.Errorf(errors.ErrInvalidLogin)
	}

	// set auth tokens
	at, rt, err := GenerateTokenPair(user)
	if err != nil {
		log.Printf("Error generating token pair for uid: %v. Error: %v\n", user.Id, err.Error())
		return dao.UserDAO{}, "", "", fmt.Errorf("failed to generate authentication tokens")
	}

	return *user, at, rt, nil
}

func Logout(userId primitive.ObjectID) error {
	token := &dao.TokenDAO{
		UserId: userId,
	}

	err := token.Delete()
	if err != nil {
		log.Printf("Error logging out user with userId: %v. Error: %v\n", userId, err.Error())
		return fmt.Errorf("failed to log user out")
	}
	return nil
}
