package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
)

type userService struct {
	userRepository  interfaces.UserRepositoryInterface
	tokenRepository interfaces.TokenRepositoryInterface
}

// NewUserService returns an interface for the user service methods
func NewUserService(userRepo interfaces.UserRepositoryInterface, tokenRepo interfaces.TokenRepositoryInterface) interfaces.UserServiceInterface {
	return &userService{
		userRepository:  userRepo,
		tokenRepository: tokenRepo,
	}
}

// Signup handles the user creation and logs the user in
func (us *userService) Signup(ctx context.Context, user *dao.User, password dto.Password) (primitive.ObjectID, error) {
	// hash the password to hide its real value
	hashedPassword, err := password.Hash()
	if err != nil {
		log.Printf("Error hashing user password: %s. Error: %v\n", password, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to sign up user", err)
	}

	// update the user password to its hashed value
	user.Password = hashedPassword

	// check that the email is not taken
	userExists, err := us.userRepository.FindByEmail(ctx, user)
	if err != nil {
		log.Printf("Error finding user with email: %s. Error: %v\n", user.Email, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to fetch user details", err)
	}

	// if the email already exists, return an error saying the email is taken
	if userExists {
		return primitive.ObjectID{}, errors.ErrBadRequest("sorry, email is taken", nil)
	}

	// create a new user with the credentials
	insertedId, err := us.userRepository.Create(ctx, user)
	if err != nil {
		log.Printf("Error creating user with email: %s. Error: %v\n", user.Email, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to sign up user", err)
	}

	return insertedId, nil
}

// Login logs the user into the application and returns the authentication tokens
func (us *userService) Login(ctx context.Context, user *dao.User, password dto.Password) error {
	// find the user by email and password
	userExists, err := us.userRepository.FindByEmail(ctx, user)
	if err != nil { // if an unexpected error occurs
		log.Printf("Error finding user with email: %s. Error: %v\n", user.Email, err.Error())
		return errors.ErrInternalServerError("failed to fetch user details", err)
	}

	// if the user does not exist, then the password and/or email are wrong
	if !userExists {
		return errors.ErrUnauthorized(errors.ErrInvalidLogin, nil)
	}

	// if the user exists, but the password is not correct
	if !password.IsEqualHash(user.Password) {
		return errors.ErrUnauthorized(errors.ErrInvalidLogin, nil)
	}

	return nil
}

// Logout logs the user out of the application
func (us *userService) Logout(ctx context.Context, userId primitive.ObjectID) error {
	token := &dao.Token{
		UserId: userId,
	}

	err := us.tokenRepository.Delete(ctx, token.UserId)
	if err != nil {
		log.Printf("Error trying to delete token with userId: %v. Error: %v\n", userId, err.Error())
		return errors.ErrInternalServerError("failed to log user out", err)
	}

	return nil
}

// GetUserByID gets a user by their ID
func (us *userService) GetUserByID(ctx context.Context, userId primitive.ObjectID) (*dao.User, error) {
	// check that the user id is not empty
	if userId.IsZero() {
		return nil, errors.ErrBadRequest("invalid user id", nil)
	}

	// create a new user object
	user := &dao.User{Id: userId}

	// find the user by the objectId
	userExists, err := us.userRepository.FindByID(ctx, user)
	if err != nil {
		log.Printf("Error finding user with id: %s. Error: %v\n", user.Id, err.Error())
		return nil, errors.ErrInternalServerError("failed to retrieve user", nil)
	}

	// return if the user does not exist
	if !userExists {
		return nil, errors.ErrBadRequest("user not found", nil)
	}

	return user, nil
}
