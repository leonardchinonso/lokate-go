package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/lokate-go/models/dao"
)

// TokenRepositoryInterface defines methods that are applicable to the token repository
type TokenRepositoryInterface interface {
	Upsert(ctx context.Context, token *dao.TokenDAO) error
	Delete(ctx context.Context, userId primitive.ObjectID) error
}

// TokenServiceInterface defines methods that are applicable to the token service
type TokenServiceInterface interface {
	GenerateTokenPair(ctx context.Context, user *dao.UserDAO) (string, string, error)
	UserFromAccessToken(tokenString string) (*dao.UserDAO, error)
}
