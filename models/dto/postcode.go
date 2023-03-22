package dto

import (
	"fmt"
	"strings"
)

// Postcode represents a UK postcode
type Postcode struct {
	StrVal  string
	byteVal []byte
	Len     int
}

// ValidatePostcode validates a postcode
func ValidatePostcode(str string) (*Postcode, error) {
	// clean up the postcode string format
	str = strings.ToUpper(strings.Replace(str, " ", "", -1))

	// begin validation
	if len(str) < 5 || len(str) > 7 {
		return nil, fmt.Errorf("invalid postcode: %s", str)
	}

	var p = &Postcode{
		StrVal: str,
		Len:    len(str),
	}
	p.byteVal = []byte(p.StrVal)

	// check that last three characters of postcode conform to UK postcode
	lastThree := p.byteVal[p.Len-3:]

	// check that first character of the last three is numeric
	if lastThree[0] < 48 || lastThree[0] > 57 {
		return nil, fmt.Errorf("invalid postcode: %s", str)
	}

	// check that the last two characters are non-numeric
	if lastThree[1] < 65 || lastThree[1] > 122 {
		return nil, fmt.Errorf("invalid postcode: %s", str)
	}
	if lastThree[2] < 65 || lastThree[2] > 122 {
		return nil, fmt.Errorf("invalid postcode: %s", str)
	}

	return p, nil
}

// NewPostcode validates and returns a Postcode object
func NewPostcode(str string) (Postcode, error) {
	p, err := ValidatePostcode(str)
	if err != nil {
		return Postcode{}, err
	}
	return *p, nil
}
