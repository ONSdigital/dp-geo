package geo

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var polygon = "Polygon"

const (
	twoPi                  float64 = 2 * math.Pi
	radiusOfEarth          float64 = 6378137 //(in metres defined by wgs84)
	defaultConcurrency             = 10      // limit number of go routines to not put too much on heap
	defaultMaximumSegments         = 180
)

// Config object to define geo configurations
type Config struct {
	defaultConcurrencyLimit int
	defaultMaxSegments      int
}

// Default geo configuration for methods on config receiver
var DefaultConfig = &Config{
	defaultConcurrencyLimit: defaultConcurrency,
	defaultMaxSegments:      defaultMaximumSegments,
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

// New instantiates a geo object with defined values to be used by geo methods
// e.g. concurrency limit and maximum number of segments when generating a circle from a point
func New(concurrencyLimit, maximumSegments int) (geo *Config) {
	geo = DefaultConfig

	geo.defaultConcurrencyLimit = concurrencyLimit
	geo.defaultMaxSegments = maximumSegments

	return
}

// NewGeoPoint creates a coordinate object
func CreateGeoPoint(lat, lon float64) (*Coordinate, error) {
	geoPoint := &Coordinate{Lat: lat, Lon: lon}

	if err := geoPoint.validate(); err != nil {
		return nil, err
	}

	return geoPoint, nil
}

// CircleToPolygon generates a bounding box (circle) from a single geopoint using radius (in metres) and the number of segments
func (geo *Config) CircleToPolygon(geoPoint Coordinate, radius float64, segments int) (*GeoStructure, error) {

	// validate input
	if err := geo.validateInput(&geoPoint, radius, segments); err != nil {
		return nil, err
	}

	shape := &GeoStructure{
		Type: polygon,
	}

	coordinates := make([][]float64, segments)

	var semaphoreChan = make(chan struct{}, geo.defaultConcurrencyLimit)

	var wg sync.WaitGroup

	for i := 0; i < segments; i++ {

		semaphoreChan <- struct{}{}

		wg.Add(1)

		go func(i int) {
			defer func() {
				<-semaphoreChan // read to release a slot
				wg.Done()
			}()

			sector := (twoPi * float64(-i)) / float64(segments)
			coordinate := generateCoordinate(geoPoint, radius, sector)
			coordinates[i] = coordinate
		}(i)
	}

	wg.Wait()

	// Push first coordinate to be last coordinate to complete polygon circle
	coordinates = append(coordinates, coordinates[0])

	shape.Coordinates = coordinates

	return shape, nil
}

func toRadians(angleInDegrees float64) float64 {
	return (angleInDegrees * math.Pi) / 180
}

func toDegrees(angleInRadians float64) float64 {
	return (angleInRadians * 180) / math.Pi
}

func generateCoordinate(geoPoint Coordinate, distance float64, sector float64) []float64 {
	lat1 := toRadians(geoPoint.Lat)
	lon1 := toRadians(geoPoint.Lon)

	// distance divided by radius of the earth
	dByR := distance / radiusOfEarth

	lat := math.Asin(
		math.Sin(lat1)*math.Cos(dByR) + math.Cos(lat1)*math.Sin(dByR)*math.Cos(sector),
	)

	lon := lon1 + math.Atan2(
		math.Sin(sector)*math.Sin(dByR)*math.Cos(lat1),
		math.Cos(dByR)-math.Sin(lat1)*math.Sin(lat),
	)

	return []float64{toDegrees(lon), toDegrees(lat)}
}

func (geo *Config) validateInput(geoPoint *Coordinate, radius float64, segments int) error {

	if err := validateRadius(radius); err != nil {
		return err
	}

	if err := validateSegments(geo.defaultMaxSegments, segments); err != nil {
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

func validateSegments(maxSegments, segments int) error {
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
