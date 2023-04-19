package interfaces

import (
	"context"

	"github.com/leonardchinonso/lokate-go/models/dao"
)

// PlaceRepositoryInterface defines methods that are applicable to the place repository
type PlaceRepositoryInterface interface {
	Create(ctx context.Context, place *dao.Place) error
	FindByID(ctx context.Context, place *dao.Place) (bool, error)
	FindByKey(ctx context.Context, place *dao.Place) (bool, error)
	PopulatePlacesInLastVisited(ctx context.Context, lastVisitedPlaces *[]dao.LastVisitedPlace) (bool, error)
	PopulatePlacesInSavedPlaces(ctx context.Context, savedPlaces *[]dao.SavedPlace) (bool, error)
}

// PlaceServiceInterface defines methods that are applicable to the place service
type PlaceServiceInterface interface {
	Create(ctx context.Context, place *dao.Place) error
	GetPlace(ctx context.Context, place *dao.Place) error
}
