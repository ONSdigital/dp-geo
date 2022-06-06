package geo

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var polygon = "Polygon"

const (
	twoPi         float64 = 2 * math.Pi
	radiusOfEarth float64 = 6378137 //(in metres defined by wgs84)
)

type maxSegments struct {
	mu       sync.Mutex
	segments int
}

var defaultMaxSegments = maxSegments{segments: 180}

// GeoStructure describes the shape of the geographical location
type GeoStructure struct {
	Type        string
	Coordinates [][]float64
}

// Coordinate describes the position of a point
type Coordinate struct {
	Lat float64
	Lon float64
}

// List of errors
var (
	ErrRadiusLargerThanEarth    = errors.New("radius can not be greater than the radius of the Earth, 6378137 metres")
	ErrTooFewSegments           = errors.New("too few segments, this should be set to 3 or more")
	ErrInvalidLongitudinalPoint = errors.New("longitude has to be between -180 and 180")
	ErrInvalidLatitudinalPoint  = errors.New("latitude has to be between -90 and 90")
)

func ErrTooManySegments(maximumSegments int) error {
	return fmt.Errorf("too many segments, this should be less than %d", maximumSegments)
}

// NewGeoPoint creates a coordinate object
func NewGeoPoint(lat, lon float64) (*Coordinate, error) {
	geoPoint := &Coordinate{Lat: lat, Lon: lon}

	if err := geoPoint.validate(); err != nil {
		return nil, err
	}

	return geoPoint, nil
}

// CircleToPolygon generates a bounding box (circle) from a single geopoint using radius (in metres) and the number of segments
func (geoPoint *Coordinate) CircleToPolygon(radius float64, segments int) (*GeoStructure, error) {

	// validate input
	if err := validateInput(geoPoint, radius, segments); err != nil {
		return nil, err
	}

	shape := &GeoStructure{
		Type: polygon,
	}

	var coordinates [][]float64

	for i := 0; i < segments; i++ {
		segment := (twoPi * float64(-i)) / float64(segments)
		coordinate := generateCoordinate(*geoPoint, radius, segment)
		coordinates = append(coordinates, coordinate)
	}

	// Push first coordinate to be last coordinate to complete polygon circle
	coordinates = append(coordinates, coordinates[0])

	shape.Coordinates = coordinates

	return shape, nil
}

// SetMaximumSegments sets the default maximum number of segments used to
// calculate the number of geo points in Polygon
func SetMaximumSegments(maxSegments int) {
	defaultMaxSegments.mu.Lock()
	defaultMaxSegments.segments = maxSegments
	defaultMaxSegments.mu.Unlock()
}

// SetMaximumSegments sets the default maximum number of segments used to
// calculate the number of geo points in Polygon
func GetMaximumSegments() (maxSegments int) {
	defaultMaxSegments.mu.Lock()
	maxSegments = defaultMaxSegments.segments
	defaultMaxSegments.mu.Unlock()

	return
}

func toRadians(angleInDegrees float64) float64 {
	return (angleInDegrees * math.Pi) / 180
}

func toDegrees(angleInRadians float64) float64 {
	return (angleInRadians * 180) / math.Pi
}

func generateCoordinate(geoPoint Coordinate, distance float64, segment float64) []float64 {
	lat1 := toRadians(geoPoint.Lat)
	lon1 := toRadians(geoPoint.Lon)

	// distance divided by radius of the earth
	dByR := distance / radiusOfEarth

	lat := math.Asin(
		math.Sin(lat1)*math.Cos(dByR) + math.Cos(lat1)*math.Sin(dByR)*math.Cos(segment),
	)

	lon := lon1 + math.Atan2(
		math.Sin(segment)*math.Sin(dByR)*math.Cos(lat1),
		math.Cos(dByR)-math.Sin(lat1)*math.Sin(lat),
	)

	return []float64{toDegrees(lon), toDegrees(lat)}
}

func validateInput(geoPoint *Coordinate, radius float64, segments int) error {

	if err := validateRadius(radius); err != nil {
		return err
	}

	if err := validateSegments(segments); err != nil {
		return err
	}

	if err := geoPoint.validate(); err != nil {
		return err
	}

	return nil
}

func validateRadius(radius float64) error {
	if radius > radiusOfEarth {
		return ErrRadiusLargerThanEarth
	}

	return nil
}

func validateSegments(segments int) error {
	maxSegments := GetMaximumSegments()
	if segments > maxSegments {
		return ErrTooManySegments(maxSegments)
	}

	if segments < 3 {
		return ErrTooFewSegments
	}

	return nil
}

func (geoPoint *Coordinate) validate() error {
	if geoPoint.Lon > 180 || geoPoint.Lon < -180 {
		return ErrInvalidLongitudinalPoint
	}

	if geoPoint.Lat > 90 || geoPoint.Lat < -90 {
		return ErrInvalidLatitudinalPoint
	}

	return nil
}
