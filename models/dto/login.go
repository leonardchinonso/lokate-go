package dto

import (
	"fmt"

	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/utils"
)

// LoginRequest holds the data for the login information
type LoginRequest struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

// Validate validates an incoming login request
func (lr *LoginRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(string(lr.Email), "email", &errs)
	utils.ShouldBePresentString(string(lr.Password), "password", &errs)

	if err := lr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	return errs
}

// LoginResponse holds the data for login response
type LoginResponse struct {
	User         dao.User `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

// NewLoginResponse returns a new LoginResponse
func NewLoginResponse(user dao.User, accessToken, refreshToken string) *LoginResponse {
	return &LoginResponse{
		User: dao.User{
			Id:          user.Id,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
