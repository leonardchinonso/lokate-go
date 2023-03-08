package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/dto"
)

// UserRepositoryInterface defines methods that are associated with the user repository
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *dao.UserDAO) (primitive.ObjectID, error)
	FindByID(ctx context.Context, user *dao.UserDAO) (bool, error)
	FindByEmail(ctx context.Context, user *dao.UserDAO) (bool, error)
}

// UserServiceInterface defines methods that are associated with the user repository
type UserServiceInterface interface {
	Signup(ctx context.Context, user *dao.UserDAO, password dto.Password) (primitive.ObjectID, error)
	Login(ctx context.Context, user *dao.UserDAO, password dto.Password) error
	Logout(ctx context.Context, userId primitive.ObjectID) error
	GetUserByID(ctx context.Context, userId primitive.ObjectID) (*dao.UserDAO, error)
}
