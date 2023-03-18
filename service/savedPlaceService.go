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
	placeRepository      interfaces.PlaceRepositoryInterface
	savedPlaceRepository interfaces.SavedPlaceRepositoryInterface
}

// NewSavedPlaceService returns an interface for the savedPlace service methods
func NewSavedPlaceService(placeRepo interfaces.PlaceRepositoryInterface, savedPlaceRepo interfaces.SavedPlaceRepositoryInterface) interfaces.SavedPlaceServiceInterface {
	return &savedPlaceService{
		placeRepository:      placeRepo,
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

	// check that the place exists
	placeExists, err := ps.placeRepository.FindByID(ctx, &dao.Place{Id: savedPlace.PlaceId})
	if err != nil {
		log.Printf("Error finding saved place by userId and placeId. Error: %v", err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	if !placeExists {
		log.Printf("Failed to retrieve place. Place does not exist")
		return errors.ErrBadRequest("cannot find place", nil)
	}

	// check that a saved place with the placeId and userId does not exist
	savedPlaceExists, err := ps.savedPlaceRepository.FindOneByPlaceIDAndUserID(ctx, savedPlace)
	if err != nil {
		log.Printf("Error finding saved place by userId and placeId. Error: %v", err)
		return errors.ErrInternalServerError(err.Error(), nil)
	}

	// if a place has already been saved by the current user
	if savedPlaceExists {
		log.Printf("Failed to saved a place. Saved place already exists for this user")
		return errors.ErrBadRequest("sorry, you have saved this place already", nil)
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

	// validate the id field
	if savedPlace.Id.IsZero() {
		log.Printf("Error validating place Id: %v\n", savedPlace.Id)
		return errors.ErrBadRequest("invalid saved_place id", nil)
	}

	// find the place by the user id and the saved place id
	placeExists, err := ps.savedPlaceRepository.FindOneByIDAndUserID(ctx, savedPlace)
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
	if savedPlace.Id.IsZero() {
		log.Printf("Error validating saved_place Id: %v\n", savedPlace.PlaceId)
		return errors.ErrBadRequest("invalid place id", nil)
	}

	// validate the savedPlace object for the PlaceAlias
	if savedPlace.PlaceAlias == "" { // if the alias is empty, make it the default of None
		savedPlace.PlaceAlias = dao.None
	} else { // else validate that it is a valid alias
		err := savedPlace.PlaceAlias.Validate()
		if err != nil {
			log.Printf("Failed to validate PlaceAlias: %v\n", savedPlace.PlaceAlias)
			return errors.ErrBadRequest("invalid place alias", nil)
		}
	}

	// reset the place alias for the user if it is a HOME or WORK
	err := ps.resetPlaceAliasForUser(ctx, savedPlace)
	if err != nil {
		log.Printf("Failed to reset PlaceAlias. Error: %v\n", err)
		return errors.ErrInternalServerError("failed to update saved place", nil)
	}

	// updated the place with the new information
	err = ps.savedPlaceRepository.Update(ctx, savedPlace)
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

// resetPlaceAliasForUser resets a user's HOME  or WORK to "NONE" to allow the new update
func (ps *savedPlaceService) resetPlaceAliasForUser(ctx context.Context, savedPlace *dao.SavedPlace) error {
	if savedPlace.PlaceAlias.IsNone() {
		return nil
	}

	// The PlaceAlias is HOME or WORK, change the former user alias to NONE
	err := ps.savedPlaceRepository.SetAlias(ctx, savedPlace, dao.None)
	if err != nil {
		log.Printf("Error resetting savedPlace PlaceAlias with id: %v and userId: %v. Error: %v\n", savedPlace.Id, savedPlace.UserId, err.Error())
		return errors.ErrInternalServerError("failed to update saved place", nil)
	}

	return nil
}
