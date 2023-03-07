package errors

import (
	"net/http"
)

const (
	ErrEmailTaken   = "email is taken"
	ErrInvalidLogin = "invalid login credentials"
)

// RestError is the custom struct for a request error
type RestError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

// ErrBadRequest returns a RestError for a bad request
func ErrBadRequest(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   "Bad Request",
		Data:    data,
	}
}

// ErrInternalServerError returns a RestError for internal server error
func ErrInternalServerError(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   "Internal Server Error",
		Data:    data,
	}
}

// ErrUnauthorized returns a RestError for an unauthorized request
func ErrUnauthorized(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusUnauthorized,
		Message: message,
		Error:   "Unauthorized",
		Data:    data,
	}
}

// ErrorToStringSlice converts a slice of errors to a slice of string
func ErrorToStringSlice(errs []error) []string {
	var errStrings []string
	for _, e := range errs {
		errStrings = append(errStrings, e.Error())
	}
	return errStrings
}
