package service

import (
	"context"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type lastVisitedPlaceService struct {
	lastVisitedPlaceRepository interfaces.LastVisitedPlaceRepositoryInterface
	placeRepository            interfaces.PlaceRepositoryInterface
}

// NewLastVisitedPlaceService returns an interface for the last visited place service
func NewLastVisitedPlaceService(lastVisitedPlaceRepo interfaces.LastVisitedPlaceRepositoryInterface, placeRepo interfaces.PlaceRepositoryInterface) interfaces.LastVisitedPlaceServiceInterface {
	return &lastVisitedPlaceService{
		lastVisitedPlaceRepository: lastVisitedPlaceRepo,
		placeRepository:            placeRepo,
	}
}

// AddLastVisitedPlace adds a place to the list of user visited places
func (l *lastVisitedPlaceService) AddLastVisitedPlace(ctx context.Context, lastVisitedPlace *dao.LastVisitedPlace) error {
	// validate the user id
	if lastVisitedPlace.UserId.IsZero() {
		log.Printf("Error validating user Id: %v\n", lastVisitedPlace.UserId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// validate the place id
	if lastVisitedPlace.PlaceId.IsZero() {
		log.Printf("Error validating place Id: %v\n", lastVisitedPlace.PlaceId)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// check that the place exists
	placeExists, err := l.placeRepository.FindByID(ctx, &dao.Place{Id: lastVisitedPlace.PlaceId})
	if err != nil {
		log.Printf("Error finding place with id: %v. Error: %v\n", lastVisitedPlace.PlaceId, err.Error())
		return errors.ErrInternalServerError("failed to retrieve place", nil)
	}

	// return an error if the place does not exist
	if !placeExists {
		return errors.ErrBadRequest("place not found", nil)
	}

	// save the details to the database
	err = l.lastVisitedPlaceRepository.Create(ctx, lastVisitedPlace)
	if err != nil {
		log.Printf("Error creating a last visited place in the database with userId: %v and placeId: %v. Error: %v",
			lastVisitedPlace.UserId, lastVisitedPlace.PlaceId, err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	return nil
}

// GetLastNVisitedPlaces gets N last visited places by a user
func (l *lastVisitedPlaceService) GetLastNVisitedPlaces(ctx context.Context, userId primitive.ObjectID, lastVisitedPlaces *[]dao.LastVisitedPlace, N int64) error {
	// validate the user id
	if userId.IsZero() {
		log.Printf("Error validating user Id: %v\n", userId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// make the database call to get the last visited N places
	placesExist, err := l.lastVisitedPlaceRepository.FindLastNVisitedPlaces(ctx, userId, lastVisitedPlaces, N)
	if err != nil {
		log.Printf("Error finding places with userId: %v. Error: %v\n", userId, err.Error())
		return errors.ErrInternalServerError("failed to retrieve places", nil)
	}

	// return an error if the places do not exist
	if !placesExist {
		return errors.ErrBadRequest("places not found", nil)
	}

	return nil
}
