package interfaces

import (
	"context"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SavedPlaceRepositoryInterface defines methods that are applicable to the savedPlace repository
type SavedPlaceRepositoryInterface interface {
	Create(ctx context.Context, savedPlace *dao.SavedPlace) error
	FindByID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error)
	FindOneByIDAndUserID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error)
	FindOneByPlaceIDAndUserID(ctx context.Context, savedPlace *dao.SavedPlace) (bool, error)
	Find(ctx context.Context, userId primitive.ObjectID, savedPlaces *[]dao.SavedPlace) (bool, error)
	Update(ctx context.Context, savedPlace *dao.SavedPlace) error
	SetAlias(ctx context.Context, savedPlace *dao.SavedPlace, newAlias dao.PlaceAlias) error
	Delete(ctx context.Context, savedPlace *dao.SavedPlace) error
}

// SavedPlaceServiceInterface defines methods that are applicable to the savedPlace service
type SavedPlaceServiceInterface interface {
	AddSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error
	GetSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error
	GetSavedPlaces(ctx context.Context, userId primitive.ObjectID, savedPlaces *[]dao.SavedPlace) error
	EditSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error
	DeleteSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error
}
