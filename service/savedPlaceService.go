package service

import (
	"context"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type savedPlaceService struct {
	savedPlaceRepository interfaces.SavedPlaceRepositoryInterface
}

// NewSavedPlaceService returns an interface for the savedPlace service methods
func NewSavedPlaceService(savedPlaceRepo interfaces.SavedPlaceRepositoryInterface) interfaces.SavedPlaceServiceInterface {
	return &savedPlaceService{
		savedPlaceRepository: savedPlaceRepo,
	}
}

// AddSavedPlace adds a savedPlace to the current user's savedPlaces
func (ps *savedPlaceService) AddSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error {
	// validate the user id field
	if savedPlace.UserId.IsZero() {
		log.Printf("Error validating user Id: %v\n", savedPlace.UserId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// validate the place id field
	if savedPlace.PlaceId.IsZero() {
		log.Printf("Error validating place Id: %v\n", savedPlace.PlaceId)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// validate the savedPlace fields
	err := savedPlace.Validate()
	if err != nil {
		log.Printf("Error validating the saved place object. Error: %v\n", err)
		return errors.ErrBadRequest(err.Error(), nil)
	}

	// save the savedPlace to the database
	err = ps.savedPlaceRepository.Create(ctx, savedPlace)
	if err != nil {
		log.Printf("Error creating a saved place with user Id: %v. Error: %v\n", savedPlace.UserId, err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	return nil
}

// GetSavedPlace gets a saved place by its id
func (ps *savedPlaceService) GetSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error {
	// validate the user id field
	if savedPlace.UserId.IsZero() {
		log.Printf("Error validating user Id: %v\n", savedPlace.UserId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// validate the place id field
	if savedPlace.PlaceId.IsZero() {
		log.Printf("Error validating place Id: %v\n", savedPlace.PlaceId)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// find the place by the user id and the place id
	placeExists, err := ps.savedPlaceRepository.FindOneByUserID(ctx, savedPlace)
	if err != nil {
		log.Printf("Error finding place with id: %s. Error: %v\n", savedPlace.Id, err.Error())
		return errors.ErrInternalServerError("failed to retrieve saved place", nil)
	}

	// return if the saved place does not exist
	if !placeExists {
		return errors.ErrBadRequest("saved place not found", nil)
	}

	return nil
}

// GetSavedPlaces gets all saved places associated with the user
func (ps *savedPlaceService) GetSavedPlaces(ctx context.Context, userId primitive.ObjectID, savedPlaces *[]dao.SavedPlace) error {
	// validate the user id field
	if userId.IsZero() {
		log.Printf("Error validating user Id: %v\n", userId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// find the places by the user id
	placesExist, err := ps.savedPlaceRepository.Find(ctx, userId, savedPlaces)
	if err != nil {
		log.Printf("Error finding places with userId: %s. Error: %v\n", userId, err.Error())
		return errors.ErrInternalServerError("failed to retrieve saved places", nil)
	}

	// return if the saved place does not exist
	if !placesExist {
		return errors.ErrBadRequest("place not found", nil)
	}

	return nil
}

// EditSavedPlace edits a SavedPlace
func (ps *savedPlaceService) EditSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error {
	// validate the user id field
	if savedPlace.UserId.IsZero() {
		log.Printf("Error validating user Id: %v\n", savedPlace.UserId)
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// validate the place id field
	if savedPlace.PlaceId.IsZero() {
		log.Printf("Error validating place Id: %v\n", savedPlace.PlaceId)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// updated the place with the new information
	err := ps.savedPlaceRepository.Update(ctx, savedPlace)
	if err != nil {
		log.Printf("Error updating place with id: %v and userId: %v. Error: %v\n", savedPlace.Id, savedPlace.UserId, err.Error())
		return errors.ErrInternalServerError("failed to update saved place", nil)
	}

	return nil
}

// DeleteSavedPlace deletes a SavedPlace to the current user's SavedPlaces
func (ps *savedPlaceService) DeleteSavedPlace(ctx context.Context, savedPlace *dao.SavedPlace) error {
	// validate the savedPlace id
	if savedPlace.Id.IsZero() {
		log.Printf("Error validating savedPlace id: %v\n", savedPlace.Id)
		return errors.ErrBadRequest("invalid savedPlace id", nil)
	}

	// delete savedPlace from database
	err := ps.savedPlaceRepository.Delete(ctx, savedPlace)
	if err != nil {
		log.Printf("Error deleting savedPlace with Id: %v. Error: %v\n", savedPlace.Id, err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	return nil
}
