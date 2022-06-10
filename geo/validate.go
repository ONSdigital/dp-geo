package geo

// validateRadius checks the input radius (in metres) is less than the radius of Earth (in metres)
func validateRadius(radius float64) error {
	if radius > radiusOfEarth {
		return ErrRadiusLargerThanEarth
	}

	return nil
}

// validateSegments checks the number of segments is not more than the maximum allowed segments
// and is not less than 3 segments
func validateSegments(maxSegments, segments int) error {
	if segments > maxSegments {
		return ErrTooManySegments(maxSegments)
	}

	if segments < 3 {
		return ErrTooFewSegments
	}

	return nil
}
