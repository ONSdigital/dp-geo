package geo

import "math"

const (
	twoPi         float64 = 2 * math.Pi
	radiusOfEarth float64 = 6378137 //(in metres defined by wgs84)
)

var polygon = "Polygon"

// GeoStructure describes the shape of the geographical location
type GeoStructure struct {
	Type        string
	Coordinates [][]float64
}

// toLowest takes in 2 values and returns the lowest value
func toLowest(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// toRadians converts angle in degrees to radians
func toRadians(angleInDegrees float64) float64 {
	return (angleInDegrees * math.Pi) / 180
}

// toDegrees converts angle in radians to degrees
func toDegrees(angleInRadians float64) float64 {
	return (angleInRadians * 180) / math.Pi
}
