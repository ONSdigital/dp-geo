package geo

import (
	"errors"
	"fmt"
)

// List of errors
var (
	ErrRadiusLargerThanEarth    = errors.New("radius can not be greater than the radius of the Earth, 6378137 metres")
	ErrTooFewSegments           = errors.New("too few segments, this should be set to 3 or more")
	ErrInvalidLongitudinalPoint = errors.New("longitude has to be between -180 and 180")
	ErrInvalidLatitudinalPoint  = errors.New("latitude has to be between -90 and 90")
)

// ErrTooManySegments returns an error containing the maximum segments allowed in error message
func ErrTooManySegments(maximumSegments int) error {
	return fmt.Errorf("too many segments, this should be less than %d", maximumSegments)
}
