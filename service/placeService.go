package service

import (
	"context"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"log"
)

type placeService struct {
	placeRepository interfaces.PlaceRepositoryInterface
}

// NewPlaceService returns an interface for the place service methods
func NewPlaceService(placeRepo interfaces.PlaceRepositoryInterface) interfaces.PlaceServiceInterface {
	return &placeService{
		placeRepository: placeRepo,
	}
}

// Create adds a place to the application
func (ps *placeService) Create(ctx context.Context, place *dao.Place) error {
	// save the place to the database
	err := ps.placeRepository.Create(ctx, place)
	if err != nil {
		log.Printf("Error creating a place. Error: %v\n", err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	return nil
}

// GetPlace gets a place by the id
func (ps *placeService) GetPlace(ctx context.Context, place *dao.Place) error {
	// validate the place id
	if place.Id.IsZero() {
		log.Printf("Error validating place id: %v\n", place.Id)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// find place in the database
	placeExists, err := ps.placeRepository.FindByID(ctx, place)
	if err != nil {
		log.Printf("Error finding place with id: %v. Error: %v\n", place.Id, err.Error())
		return errors.ErrInternalServerError("failed to retrieve place", nil)
	}

	// return an error if the place does not exist
	if !placeExists {
		return errors.ErrBadRequest("place not found", nil)
	}

	return nil
}
