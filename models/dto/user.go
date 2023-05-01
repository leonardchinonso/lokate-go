package dto

import (
	"fmt"

	"github.com/leonardchinonso/lokate-go/utils"
)

// EditProfileRequest holds the data for the edit profile information
type EditProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       Email  `json:"email"`
}

// Validate validates an incoming edit profile request
func (epr *EditProfileRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(epr.FirstName, "first_name", &errs)
	utils.ShouldBePresentString(epr.LastName, "last_name", &errs)
	utils.ShouldBePresentString(string(epr.Email), "email", &errs)

	// validate the email
	if err := epr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	return errs
}
