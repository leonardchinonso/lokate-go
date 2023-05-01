package dto

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/utils"
)

// SignupRequest holds the data for signup information
type SignupRequest struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Email           Email    `json:"email"`
	Password        Password `json:"password"`
	ConfirmPassword Password `json:"confirm_password"`
}

// Validate validates an incoming signup request
func (sr *SignupRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(sr.FirstName, "first name", &errs)
	utils.ShouldBePresentString(sr.LastName, "last name", &errs)
	utils.ShouldBePresentString(string(sr.Email), "email", &errs)
	utils.ShouldBePresentString(string(sr.Password), "password", &errs)
	utils.ShouldBePresentString(string(sr.ConfirmPassword), "confirmed password", &errs)

	// validate the email
	if err := sr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	// validate the password
	if err := sr.Password.Validate(); err != nil {
		errs = append(errs, err)
	} else if ok := sr.Password.IsEqualValue(sr.ConfirmPassword); !ok {
		errs = append(errs, fmt.Errorf("passwords do not match"))
	}

	return errs
}
