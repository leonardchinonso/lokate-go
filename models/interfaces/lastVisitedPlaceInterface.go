package interfaces

import (
	"context"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LastVisitedPlaceRepositoryInterface holds the methods for accessing the last visited place repository
type LastVisitedPlaceRepositoryInterface interface {
	Create(ctx context.Context, lastVisitedPlace *dao.LastVisitedPlace) error
	FindLastNVisitedPlaces(ctx context.Context, UserId primitive.ObjectID, lastVisitedPlace *[]dao.LastVisitedPlace, N int64) (bool, error)
}

// LastVisitedPlaceServiceInterface holds the methods for accessing the last visited place service
type LastVisitedPlaceServiceInterface interface {
	AddLastVisitedPlace(ctx context.Context, lastVisitedPlace *dao.LastVisitedPlace) error
	GetLastNVisitedPlaces(ctx context.Context, userId primitive.ObjectID, lastVisitedPlaces *[]dao.LastVisitedPlace, N int64) error
}
