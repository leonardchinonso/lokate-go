package utils

import (
	"fmt"
)

// ShouldBePresentString checks that a string in a field is required
func ShouldBePresentString(value string, prefix string, errs *[]error) {
	if value == "" {
		*errs = append(*errs, fmt.Errorf(fmt.Sprintf("%s is required", prefix)))
	}
}
